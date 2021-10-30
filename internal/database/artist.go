package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"go.uber.org/multierr"

	"github.com/obitech/artist-db/internal/database/model"
)

// UpsertArtists creates or updates one or more artists in the database.
// Multiple artists are inserted in the same transaction
func (db *Database) UpsertArtists(ctx context.Context, artists ...*model.Artist) error {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("creating tx failed: %w", err)
	}

	defer rollbackAndLogError(ctx, tx)

	var mErr error
	for _, artist := range artists {
		if err := db.upsertArtist(ctx, tx, artist); err != nil {
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

func (db *Database) upsertArtist(ctx context.Context, tx pgx.Tx, artist *model.Artist) error {
	stmt := fmt.Sprintf(`
		INSERT INTO "%s"
			(
				id, 
				first_name,
				last_name,
				pronouns,
				date_of_birth,
				place_of_birth,
				nationality,
				language,
				facebook,
				instagram,
				bandcamp,
				bio_ger,
				bio_en,
				artist_name
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT 
			(id)
		DO UPDATE SET
			first_name=$2,
			last_name=$3,
			pronouns=$4,
			date_of_birth=$5,
			place_of_birth=$6,
			nationality=$7,
			language=$8,
			facebook=$9,
			instagram=$10,
			bandcamp=$11,
			bio_ger=$12,
			bio_en=$13,
			artist_name=$14`, TableArtists)

	if _, err := tx.Exec(ctx, stmt,
		artist.ID,                       // $1
		artist.FirstName,                // $2
		artist.LastName,                 // $3
		artist.Pronouns,                 // $4
		artist.Origin.DateOfBirth.UTC(), // $5
		artist.Origin.PlaceOfBirth,      // $6
		artist.Origin.Nationality,       // $7
		artist.Language,                 // $8
		artist.Socials.Facebook,         // $9
		artist.Socials.Instagram,        // $10
		artist.Socials.Bandcamp,         // $11
		artist.BioGerman,                // $12
		artist.BioEnglish,               // $13
		artist.ArtistName,               // $14
	); err != nil {
		return err
	}

	return nil
}
