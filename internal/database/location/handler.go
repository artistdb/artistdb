package location

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/obitech/artist-db/internal/database/core"
)

// Handler returns a DB Handler that operates on Locations.
type Handler struct {
	conn   core.Connection
	logger *zap.Logger
}

// NewHandler returns a handler.
func NewHandler(conn core.Connection, logger *zap.Logger) *Handler {
	return &Handler{
		conn:   conn,
		logger: logger,
	}
}

// Upsert inserts or updates Locations.
func (h *Handler) Upsert(ctx context.Context, locations ...*Location) error {
	tx, err := h.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("creating tx failed: %w", err)
	}

	defer core.RollbackAndLogError(ctx, tx, h.logger)

	var mErr error
	for _, location := range locations {
		if err := h.upsert(ctx, tx, location); err != nil {
			if errors.Is(err, pgx.ErrTxClosed) {
				return fmt.Errorf("insert aborted, tx cancelled: %w", err)
			}

			mErr = multierr.Append(mErr, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commiting tx failed: %w", err)
	}

	return mErr
}

func (h *Handler) upsert(ctx context.Context, tx pgx.Tx, location *Location) error {
	start := time.Now().UTC()

	stmt := fmt.Sprintf(`
		INSERT INTO "%s"
			(
				id,
				created_at,
				updated_at,
				name
			)
		VALUES 
			($1, $2, $3, $4)
		ON CONFLICT 
			(id)
		DO UPDATE SET
			updated_at=$3,
			name=$4,
			deleted_at=NULL`, core.TableLocations)

	if _, err := tx.Exec(ctx, stmt,
		location.ID,   // $1
		start,         // $2
		start,         // $3
		location.Name, // $4
	); err != nil {
		return err
	}

	return nil
}
