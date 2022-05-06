package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	otelTrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/multierr"
	"go.uber.org/zap"

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
		if err := h.upsertEvents(spanCtx, tx, event); err != nil {
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

func (h *Handler) upsertEvents(ctx context.Context, tx pgx.Tx, event *Event) error {
	var (
		start      = time.Now().UTC()
		locationID *string
	)

	if loc := event.Location; loc != nil {
		if loc.ID == "" {
			return fmt.Errorf("location ID is empty")
		}

		locationID = &loc.ID
	}

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

	if _, err := tx.Exec(ctx, stmt,
		event.ID,
		start,
		start,
		event.Name,
		event.StartTime.UTC(),
		locationID,
	); err != nil {
		return err
	}

	return nil
}
