package artist

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/obitech/artist-db/internal/conversion"
	"github.com/obitech/artist-db/internal/database/core"
)

// Handler returns a DB Handler which operates on Artists.
type Handler struct {
	conn   core.Connection
	logger *zap.Logger
}

// NewHandler returns a Handler.
func NewHandler(conn core.Connection, logger *zap.Logger) *Handler {
	return &Handler{
		conn:   conn,
		logger: logger,
	}
}

// Upsert creates or updates one or more artists in the database.
// Multiple artists are inserted in the same transaction
func (h *Handler) Upsert(ctx context.Context, artists ...*Artist) error {
	tx, err := h.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("creating tx failed: %w", err)
	}

	defer core.RollbackAndLogError(ctx, tx, h.logger)

	var mErr error
	for _, artist := range artists {
		if err := h.upsertArtist(ctx, tx, artist); err != nil {
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

func (h *Handler) upsertArtist(ctx context.Context, tx pgx.Tx, artist *Artist) error {
	start := time.Now().UTC()

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
				artist_name,
				created_at,
				updated_at
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
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
			artist_name=$14,
			updated_at=$16,
			deleted_at=NULL`, core.TableArtists)

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
		start,                           // $15
		start,                           // $16
	); err != nil {
		return err
	}

	return nil
}

// GetRequest specifies the input for an  Artists query against the database.
type GetRequest func() (string, string)

// ByID requests and Artist by ID.
func ByID(id string) GetRequest {
	return func() (string, string) {
		return id, "id=$1"
	}
}

// ByArtistName requests Artists by the artists'.
func ByArtistName(firstName string) GetRequest {
	return func() (string, string) {
		return firstName, "artist_name=$1"
	}
}

// ByLastName requests Artists by last name.
func ByLastName(lastName string) GetRequest {
	return func() (string, string) {
		return lastName, "last_name=$1"
	}
}

// Get retrieves Artists according to GetRequest, or an ErrNotFound.
func (h *Handler) Get(ctx context.Context, request GetRequest) ([]*Artist, error) {
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
		WHERE deleted_at IS NULL AND `, core.TableArtists,
	)

	stmt += whereClause

	rows, err := h.conn.Query(ctx, stmt, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.ErrNotFound
		}

		return nil, fmt.Errorf("query failed: %w", err)
	}

	defer rows.Close()

	var artists []*Artist

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

		artists = append(artists, &Artist{
			ID:         id,
			FirstName:  firstName,
			LastName:   lastName,
			ArtistName: conversion.String(artistName),
			Pronouns:   pronouns,
			Origin: Origin{
				DateOfBirth:  conversion.Time(dob),
				PlaceOfBirth: conversion.String(pob),
				Nationality:  conversion.String(nationality),
			},
			Language: conversion.String(language),
			Socials: Socials{
				Instagram: conversion.String(instagram),
				Facebook:  conversion.String(facebook),
				Bandcamp:  conversion.String(bandcamp),
			},
			BioGerman:  conversion.String(bioGer),
			BioEnglish: conversion.String(bioEn),
		})

	}

	if len(artists) == 0 {
		return nil, core.ErrNotFound
	}

	return artists, nil
}

// DeleteByID deletes an Artist by ID. Returns ErrNotFound if the Artist
// did not exist beforehand.
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
			id`, core.TableArtists)

	var deletedID string
	if err := h.conn.QueryRow(ctx, stmt, conversion.TimeP(time.Now().UTC()), id).Scan(&deletedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core.ErrNotFound
		}

		return err
	}

	if deletedID == "" {
		return core.ErrNotFound
	}

	return nil
}