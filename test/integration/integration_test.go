package integration

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/obitech/artist-db/internal/database"
	"github.com/obitech/artist-db/internal/database/model"
)

func setup(t *testing.T, ctx context.Context) (*database.Database, *pgx.Conn, func(t *testing.T)) {
	connString := os.Getenv("TEST_DB_CONN_STRING")
	require.NotEmpty(t, connString)

	var (
		db  *database.Database
		err error
	)

	// Wait for the database to come up.
	do := func() bool {
		db, err = database.NewDatabase(ctx, connString)
		return err == nil && db.Ready(ctx) == nil
	}

	// Wait for the database to be ready. This might take some time in CI.
	require.Eventuallyf(t, do, 30*time.Second, 500*time.Millisecond, "database didn't come up")

	require.NoError(t, db.CreateTables(connString))

	conn, err := pgx.Connect(ctx, connString)
	require.NoError(t, err)

	teardown := func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		require.NoError(t, db.DestroyTables(connString))
		require.NoError(t, conn.Close(ctx))
	}

	return db, conn, teardown
}

func toString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func Test_TablesExistsIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, conn, teardown := setup(t, ctx)
	defer teardown(t)

	stmt := `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema='public' and table_name=$1);`

	t.Run("artists exists", func(t *testing.T) {
		var exists bool
		require.NoError(t, conn.QueryRow(ctx, stmt, database.TableArtists).Scan(&exists))

		assert.True(t, exists)
	})

	t.Run("artworks exists", func(t *testing.T) {
		var exists bool
		require.NoError(t, conn.QueryRow(ctx, stmt, database.TableArtworks).Scan(&exists))

		assert.True(t, exists)
	})

	t.Run("events exists", func(t *testing.T) {
		var exists bool
		require.NoError(t, conn.QueryRow(ctx, stmt, database.TableEvents).Scan(&exists))

		assert.True(t, exists)
	})

	t.Run("locations exists", func(t *testing.T) {
		var exists bool
		require.NoError(t, conn.QueryRow(ctx, stmt, database.TableLocations).Scan(&exists))

		assert.True(t, exists)
	})

	t.Run("artwork_event_locations exists", func(t *testing.T) {
		var exists bool
		require.NoError(t, conn.QueryRow(ctx, stmt, database.TableArtworkEventLocations).Scan(&exists))

		assert.True(t, exists)
	})

	t.Run("invited_artists exists", func(t *testing.T) {
		var exists bool
		require.NoError(t, conn.QueryRow(ctx, stmt, database.TableInvitedArtists).Scan(&exists))

		assert.True(t, exists)
	})
}

func Test_ArtistsIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, conn, teardown := setup(t, ctx)
	defer teardown(t)

	artists := []*model.Artist{
		{
			ID:         uuid.New().String(),
			FirstName:  "first",
			LastName:   "last",
			ArtistName: "artist",
			Pronouns:   []string{"she", "her"},
			Origin: model.Origin{
				DateOfBirth:  time.Time{},
				PlaceOfBirth: "de",
				Nationality:  "de",
			},
			Language: "de",
			Socials: model.Socials{
				Instagram: "foo",
				Facebook:  "bar",
				Bandcamp:  "baz",
			},
			BioGerman:  "alfred",
			BioEnglish: "biolek",
		},
		{
			ID:         uuid.New().String(),
			FirstName:  "first2",
			LastName:   "last2",
			ArtistName: "artist2",
			Pronouns:   []string{"he", "her", "him"},
			Origin: model.Origin{
				DateOfBirth:  time.Time{}.Add(time.Hour),
				PlaceOfBirth: "en",
				Nationality:  "en",
			},
			Language: "en",
			Socials: model.Socials{
				Instagram: "foo2",
				Facebook:  "bar2",
				Bandcamp:  "baz2",
			},
			BioGerman:  "alfred2",
			BioEnglish: "biolek2",
		},
	}

	t.Run("inserting and retrieving single artist works", func(t *testing.T) {
		require.NoError(t, db.UpsertArtists(ctx, artists[0]))

		t.Run("verify", func(t *testing.T) {
			artist, err := db.GetArtistByID(ctx, artists[0].ID)
			require.NoError(t, err)

			assert.Equal(t, artists[0], artist)
		})

		t.Run("cleanup", func(t *testing.T) {
			stmt := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, database.TableArtists)

			_, err := conn.Exec(ctx, stmt, artists[0].ID)
			require.NoError(t, err)

			stmt = fmt.Sprintf(`SELECT id from %s WHERE id=$1`, database.TableArtists)

			var id string
			require.Error(t, conn.QueryRow(ctx, stmt, artists[0].ID).Scan(&id))
			assert.Empty(t, id, "")
		})
	})

	t.Run("inserting multiple artists works", func(t *testing.T) {
		require.NoError(t, db.UpsertArtists(ctx, artists...))

		t.Run("verify", func(t *testing.T) {
			artist, err := db.GetArtistByID(ctx, artists[0].ID)
			require.NoError(t, err)
			assert.Equal(t, artists[0], artist)

			artist, err = db.GetArtistByID(ctx, artists[1].ID)
			require.NoError(t, err)
			assert.Equal(t, artists[1], artist)
		})

		t.Run("updating existing artist works", func(t *testing.T) {
			artists[0].ArtistName = "pee.age"
			require.NoError(t, db.UpsertArtists(ctx, artists...))

			artist, err := db.GetArtistByID(ctx, artists[0].ID)
			require.NoError(t, err)
			assert.Equal(t, artists[0], artist)
			assert.Equal(t, "pee.age", artist.ArtistName)
		})
	})

	t.Run("retrieving non-existent artist throws error", func(t *testing.T) {
		t.Run("invalid ID throws error", func(t *testing.T) {
			artist, err := db.GetArtistByID(ctx, "foo")
			require.Error(t, err)

			assert.True(t, errors.Is(err, database.ErrInvalidUUID), err.Error())
			assert.Nil(t, artist)
		})

		t.Run("unknown ID thwrows error", func(t *testing.T) {
			artist, err := db.GetArtistByID(ctx, uuid.New().String())
			require.Error(t, err)

			assert.True(t, errors.Is(err, database.ErrNotFound), err.Error())
			assert.Nil(t, artist)
		})
	})
}
