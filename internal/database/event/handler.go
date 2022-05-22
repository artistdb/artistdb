package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.opentelemetry.io/otel/attribute"
	otelTrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/obitech/artist-db/internal/conversion"
	"github.com/obitech/artist-db/internal/database/core"
	"github.com/obitech/artist-db/internal/observability"
)

const (
	entityEvent = "event"
)

// Handler is a DB Handler which operates on Events.
type Handler struct {
	conn   core.Connection
	logger *zap.Logger
	tracer otelTrace.TracerProvider
}

// NewHandler returns an events Handler..
func NewHandler(conn core.Connection, logger *zap.Logger, tp otelTrace.TracerProvider) *Handler {
	return &Handler{
		conn:   conn,
		logger: logger,
		tracer: tp,
	}
}

// Upsert creates or updates one or more events in the database.
// Multiple events are inserted in the same transaction.
func (h *Handler) Upsert(ctx context.Context, events ...*Event) error {
	spanCtx, span := h.tracer.Tracer(core.TracingInstrumentationName).Start(ctx, "artist.upsert")
	defer span.End()

	tx, err := h.conn.Begin(spanCtx)
	if err != nil {
		return fmt.Errorf("creating tx failed: %w", err)
	}

	defer core.RollbackAndLogError(spanCtx, tx, h.logger)

	var (
		mErr          error
		eventsChanged int
	)
	for _, event := range events {
		if err := h.upsertEvent(spanCtx, tx, event); err != nil {
			if errors.Is(err, pgx.ErrTxClosed) {
				return fmt.Errorf("insert aborted, tx cancelled: %w", err)
			}

			observability.Metrics.TrackObjectError(entityEvent, "upsert")
			mErr = multierr.Append(mErr, err)
		} else {
			eventsChanged++
		}
	}

	if err := tx.Commit(spanCtx); err != nil {
		return fmt.Errorf("commiting tx failed: %w", err)
	}

	observability.Metrics.TrackObjectsChanged(eventsChanged, entityEvent, "upsert")

	return mErr
}

func (h *Handler) upsertEvent(ctx context.Context, tx pgx.Tx, event *Event) error {
	start := time.Now().UTC()

	stmt := fmt.Sprintf(`
		INSERT INTO "%s"
			(
				id,
				created_at,
				updated_at,
				name,
				start_time,
				location_id
			)
		VALUES
			($1, $2, $3, $4, $5, $6)
		ON CONFLICT
			(id)
		DO UPDATE SET
			name=$4,
			start_time=$5,
			location_id=$6,
			deleted_at=NULL`, core.TableEvents)

	if event.StartTime != nil {
		t := *event.StartTime
		t = t.UTC()
		event.StartTime = &t
	}

	if _, err := tx.Exec(ctx, stmt,
		event.ID,
		start,
		start,
		event.Name,
		event.StartTime,
		event.LocationID,
	); err != nil {
		return fmt.Errorf("upserting event: %w", err)
	}

	for _, invitedArtist := range event.InvitedArtists {
		if err := h.inviteArtist(ctx, tx, event.ID, invitedArtist); err != nil {
			return fmt.Errorf("upsert invited artist: %w", err)
		}
	}

	return nil
}

func (h *Handler) inviteArtist(ctx context.Context, tx pgx.Tx, eventID string, invitedArtist InvitedArtist) error {
	stmt := fmt.Sprintf(`
		INSERT INTO %q
			(
				artist_id,
				event_id,
				confirmed
			)
		VALUES
			($1, $2, $3)
		ON CONFLICT
			(artist_id, event_id)
		DO UPDATE SET
			confirmed=$3`, core.TableInvitedArtists)

	if _, err := tx.Exec(ctx, stmt,
		invitedArtist.ID,
		eventID,
		invitedArtist.Confirmed,
	); err != nil {
		return err
	}

	return nil
}

// GetRequest specifies the input for an Event query against the database.
type GetRequest func() (string, string, string)

// ByID requests an Event by ID.
func ByID(id string) GetRequest {
	return func() (string, string, string) {
		return id, fmt.Sprintf("%s.id=$1", core.TableEvents), "id"
	}
}

// ByName requests an Event by name.
func ByName(name string) GetRequest {
	return func() (string, string, string) {
		return name, fmt.Sprintf("%s.name=$1", core.TableEvents), "name"
	}
}

func (h *Handler) Get(ctx context.Context, req GetRequest) ([]*Event, error) {
	input, whereClause, reqType := req()

	spanCtx, span := h.tracer.Tracer(core.TracingInstrumentationName).Start(ctx, "artist.get", otelTrace.WithAttributes(attribute.String("type", reqType)))
	defer span.End()

	stmt := fmt.Sprintf(`
		SELECT 
			id,
			name,
			start_time,
			location_id
		FROM "%s"
		WHERE 
			deleted_at IS NULL AND `, core.TableEvents) + whereClause

	rows, err := h.conn.Query(spanCtx, stmt, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.ErrNotFound
		}

		observability.Metrics.TrackObjectError(entityEvent, "get")
		return nil, fmt.Errorf("query failed: %w", err)
	}

	defer rows.Close()

	var events []*Event

	for rows.Next() {
		var (
			id         string
			name       string
			startTime  *time.Time
			locationID *string
		)

		if err := rows.Scan(&id, &name, &startTime, &locationID); err != nil {
			span.RecordError(err)
			observability.Metrics.TrackObjectError(entityEvent, "get")
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		if startTime != nil {
			var t time.Time
			t = *startTime
			t = t.UTC()
			startTime = &t
		}

		invited, err := h.invitedArtists(ctx, id)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("retrieving invtied artists: %w", err)
		}

		events = append(events, &Event{
			ID:             id,
			Name:           name,
			StartTime:      startTime,
			LocationID:     locationID,
			InvitedArtists: invited,
		})
	}

	if len(events) == 0 {
		return nil, core.ErrNotFound
	}

	observability.Metrics.TrackObjectsRetrieved(len(events), entityEvent)

	return events, nil
}

func (h *Handler) DeleteByID(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return core.ErrInvalidUUID
	}

	spanCtx, span := h.tracer.Tracer(core.TracingInstrumentationName).Start(ctx, "event.delete")
	defer span.End()

	stmt := fmt.Sprintf(`
		UPDATE 
			"%s" 
		SET 
			deleted_at=$1,
			updated_at=$1
		WHERE 
			id=$2 
		RETURNING 
			id`, core.TableEvents)

	var deletedID string
	if err := h.conn.QueryRow(spanCtx, stmt, conversion.TimeP(time.Now().UTC()), id).Scan(&deletedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core.ErrNotFound
		}

		observability.Metrics.TrackObjectError(entityEvent, "delete")
		span.RecordError(err)
		return err
	}

	if deletedID == "" {
		return core.ErrNotFound
	}

	observability.Metrics.TrackObjectsChanged(1, entityEvent, "delete")

	return nil
}

func (h *Handler) invitedArtists(ctx context.Context, eventID string) ([]InvitedArtist, error) {
	stmt := fmt.Sprintf(`
		SELECT
			artist_id, confirmed
		FROM
			%q
		WHERE
			event_id=$1`, core.TableInvitedArtists)

	rows, err := h.conn.Query(ctx, stmt, eventID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer rows.Close()

	var invited []InvitedArtist
	for rows.Next() {
		var (
			id        string
			confirmed bool
		)

		if err := rows.Scan(&id, &confirmed); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		invited = append(invited, InvitedArtist{
			ID:        id,
			Confirmed: confirmed,
		})
	}

	return invited, nil
}
