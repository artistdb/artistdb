package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"go.uber.org/multierr"

	"github.com/obitech/artist-db/internal/database/model"
)

func (db *Database) UpsertLocations(ctx context.Context, locations ...*model.Location) error {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("creating tx failed: %w", err)
	}

	defer rollbackAndLogError(ctx, tx)

	var mErr error
	for _, location := range locations {
		if err := db.upsertLocation(ctx, tx, location); err != nil {
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

func (db *Database) upsertLocation(ctx context.Context, tx pgx.Tx, location *model.Location) error {
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
			deleted_at=NULL`, TableLocations)

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
