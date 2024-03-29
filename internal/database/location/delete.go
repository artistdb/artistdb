package location

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/obitech/artist-db/internal/conversion"
	"github.com/obitech/artist-db/internal/database/core"
	"github.com/obitech/artist-db/internal/observability"
)

const (
	entityLocation = "locations"
)

// DeleteByID deletes a Location by ID.
func (h *Handler) DeleteByID(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return core.ErrInvalidUUID
	}

	stmt := fmt.Sprintf(`
		UPDATE 
			"%s" 
		SET 
			deleted_at=$1,
			updated_at=$1
		WHERE 
			id=$2 
		RETURNING 
			id`, core.TableLocations)

	var deletedID string
	if err := h.conn.QueryRow(ctx, stmt, conversion.TimeP(time.Now().UTC()), id).Scan(&deletedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core.ErrNotFound
		}

		observability.Metrics.TrackObjectError(entityLocation, "delete")
		return err
	}

	if deletedID == "" {
		return core.ErrNotFound
	}

	observability.Metrics.TrackObjectsChanged(1, entityLocation, "delete")
	h.logger.Info("tuple modified",
		zap.String("action", "delete"),
		zap.String("entity", entityLocation),
		zap.String("id", id),
	)

	return nil
}
