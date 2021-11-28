package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.uber.org/multierr"

	"github.com/obitech/artist-db/internal/conversion"
	"github.com/obitech/artist-db/internal/database/model"
)

var (
	ErrNotFound    = errors.New("resource not found")
	ErrInvalidUUID = errors.New("id must be valid UUID")
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

// GetArtistRequest specifies the input for an  Artists query against the database.
type GetArtistRequest func() (string, string)

// ByID requests and Artist by ID.
func ByID(id string) GetArtistRequest {
	return func() (string, string) {
		return id, "id=$1"
	}
}

// ByArtistName requests Artists by the artists'.
func ByArtistName(firstName string) GetArtistRequest {
	return func() (string, string) {
		return firstName, "artist_name=$1"
	}
}

// ByLastName requests Artists by last name.
func ByLastName(lastName string) GetArtistRequest {
	return func() (string, string) {
		return lastName, "last_name=$1"
	}
}

// GetArtists retrieves Artists according to GetArtistRequest, or an ErrNotFound.
func (db *Database) GetArtists(ctx context.Context, request GetArtistRequest) ([]*model.Artist, error) {
	input, whereClause := request()

	stmt := fmt.Sprintf(`
		SELECT 
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
		FROM
			"%s"
		WHERE `, TableArtists,
	)

	stmt += whereClause

	rows, err := db.conn.Query(ctx, stmt, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("query failed: %w", err)
	}

	defer rows.Close()

	var artists []*model.Artist

	for rows.Next() {
		var (
			id          string
			firstName   string
			lastName    string
			pronouns    []string
			dob         *time.Time
			pob         *string
			nationality *string
			language    *string
			facebook    *string
			instagram   *string
			bandcamp    *string
			bioGer      *string
			bioEn       *string
			artistName  *string
		)

		if err := rows.Scan(
			&id,
			&firstName,
			&lastName,
			&pronouns,
			&dob,
			&pob,
			&nationality,
			&language,
			&facebook,
			&instagram,
			&bandcamp,
			&bioGer,
			&bioEn,
			&artistName,
		); err != nil {
			return nil, fmt.Errorf("scanning rows failed: %w", err)
		}

		artists = append(artists, &model.Artist{
			ID:         id,
			FirstName:  firstName,
			LastName:   lastName,
			ArtistName: conversion.PointerToString(artistName),
			Pronouns:   pronouns,
			Origin: model.Origin{
				DateOfBirth:  conversion.PointerToTime(dob),
				PlaceOfBirth: conversion.PointerToString(pob),
				Nationality:  conversion.PointerToString(nationality),
			},
			Language: conversion.PointerToString(language),
			Socials: model.Socials{
				Instagram: conversion.PointerToString(instagram),
				Facebook:  conversion.PointerToString(facebook),
				Bandcamp:  conversion.PointerToString(bandcamp),
			},
			BioGerman:  conversion.PointerToString(bioGer),
			BioEnglish: conversion.PointerToString(bioEn),
		})

	}

	if len(artists) == 0 {
		return nil, ErrNotFound
	}

	return artists, nil
}

// DeleteArtistByID deletes an Artist by ID. Returns ErrNotFound if the Artist
// did not exist beforehand.
func (db *Database) DeleteArtistByID(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidUUID
	}

	stmt := fmt.Sprintf(`DELETE FROM "%s" WHERE id=$1 RETURNING id`, TableArtists)

	var deletedID string
	if err := db.conn.QueryRow(ctx, stmt, id).Scan(&deletedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}

		return err
	}

	if deletedID == "" {
		return ErrNotFound
	}

	return nil
}
