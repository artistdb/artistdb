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

	// t.Run("artworks exists", func(t *testing.T) {
	// 	var exists bool
	// 	require.NoError(t, conn.QueryRow(ctx, stmt, database.TableArtworks).Scan(&exists))
	//
	// 	assert.True(t, exists)
	// })
	//
	// t.Run("events exists", func(t *testing.T) {
	// 	var exists bool
	// 	require.NoError(t, conn.QueryRow(ctx, stmt, database.TableEvents).Scan(&exists))
	//
	// 	assert.True(t, exists)
	// })
	//
	// t.Run("locations exists", func(t *testing.T) {
	// 	var exists bool
	// 	require.NoError(t, conn.QueryRow(ctx, stmt, database.TableLocations).Scan(&exists))
	//
	// 	assert.True(t, exists)
	// })
	//
	// t.Run("artwork_event_locations exists", func(t *testing.T) {
	// 	var exists bool
	// 	require.NoError(t, conn.QueryRow(ctx, stmt, database.TableArtworkEventLocations).Scan(&exists))
	//
	// 	assert.True(t, exists)
	// })
	//
	// t.Run("invited_artists exists", func(t *testing.T) {
	// 	var exists bool
	// 	require.NoError(t, conn.QueryRow(ctx, stmt, database.TableInvitedArtists).Scan(&exists))
	//
	// 	assert.True(t, exists)
	// })
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
		{
			ID:         uuid.New().String(),
			FirstName:  "first2",
			LastName:   "last2",
			ArtistName: "artist2",
		},
	}

	t.Run("inserting and retrieving single artist works", func(t *testing.T) {
		t.Run("invalid ID throws error", func(t *testing.T) {
			require.Error(t, db.UpsertArtists(ctx, &model.Artist{ID: "foo"}))
		})

		require.NoError(t, db.UpsertArtists(ctx, artists[0]))

		t.Run("verify", func(t *testing.T) {
			t.Run("resources are created", func(t *testing.T) {
				res, err := db.GetArtists(ctx, database.ByID(artists[0].ID))
				require.NoError(t, err)
				require.Len(t, res, 1)
				assert.Equal(t, artists[0], res[0])

				res, err = db.GetArtists(ctx, database.ByArtistName(artists[0].ArtistName))
				require.NoError(t, err)
				require.Len(t, res, 1)
				assert.Equal(t, artists[0], res[0])

				res, err = db.GetArtists(ctx, database.ByLastName(artists[0].LastName))
				require.NoError(t, err)
				require.Len(t, res, 1)
				assert.Equal(t, artists[0], res[0])
			})

			t.Run("metadata is set", func(t *testing.T) {
				stmt := fmt.Sprintf(`SELECT created_at, updated_at, deleted_at FROM %s WHERE ID=$1`, database.TableArtists)

				var (
					createdAt time.Time
					updatedAt time.Time
					deletedAt *time.Time
				)

				require.NoError(t, conn.QueryRow(ctx, stmt, artists[0].ID).Scan(&createdAt, &updatedAt, &deletedAt))

				assert.NotZero(t, createdAt)
				assert.NotZero(t, updatedAt)
				assert.Equal(t, updatedAt, createdAt)
				assert.Nil(t, deletedAt)
			})
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
			res, err := db.GetArtists(ctx, database.ByID(artists[0].ID))
			require.NoError(t, err)
			require.Len(t, res, 1)
			assert.Equal(t, artists[0], res[0])

			res, err = db.GetArtists(ctx, database.ByID(artists[1].ID))
			require.NoError(t, err)
			require.Len(t, res, 1)
			assert.Equal(t, artists[1], res[0])

			res, err = db.GetArtists(ctx, database.ByArtistName(artists[1].ArtistName))
			require.NoError(t, err)
			require.Len(t, res, 2)
			assert.Equal(t, artists[1], res[0])
			assert.Equal(t, artists[2], res[1])

			res, err = db.GetArtists(ctx, database.ByLastName(artists[1].LastName))
			require.NoError(t, err)
			require.Len(t, res, 2)
			assert.Equal(t, artists[1], res[0])
			assert.Equal(t, artists[2], res[1])
		})

		t.Run("updating existing artist works", func(t *testing.T) {
			artists[0].ArtistName = "pee.age"
			require.NoError(t, db.UpsertArtists(ctx, artists...))

			res, err := db.GetArtists(ctx, database.ByID(artists[0].ID))
			require.NoError(t, err)
			require.Len(t, res, 1)
			assert.Equal(t, artists[0], res[0])
			assert.Equal(t, "pee.age", res[0].ArtistName)
		})
	})

	t.Run("retrieving non-existent artist throws error", func(t *testing.T) {
		t.Run("invalid ID throws error", func(t *testing.T) {
			res, err := db.GetArtists(ctx, database.ByID("foo"))
			require.Error(t, err)
			assert.Nil(t, res)
		})

		t.Run("unknown ID throws error", func(t *testing.T) {
			artist, err := db.GetArtists(ctx, database.ByID(uuid.New().String()))
			require.Error(t, err)

			assert.True(t, errors.Is(err, database.ErrNotFound), err.Error())
			assert.Nil(t, artist)
		})
	})

	t.Run("deleting non-existent artist throws error", func(t *testing.T) {
		t.Run("invalid ID throws error", func(t *testing.T) {
			err := db.DeleteArtistByID(ctx, "foo")
			require.Error(t, err)

			assert.True(t, errors.Is(err, database.ErrInvalidUUID), err.Error())
		})

		t.Run("unknown ID thwrows error", func(t *testing.T) {
			err := db.DeleteArtistByID(ctx, uuid.New().String())
			require.Error(t, err)

			assert.True(t, errors.Is(err, database.ErrNotFound), err.Error())
		})
	})

	t.Run("deleting artist works", func(t *testing.T) {
		t.Run("delete", func(t *testing.T) {
			err := db.DeleteArtistByID(ctx, artists[0].ID)
			require.NoError(t, err)
		})

		t.Run("validate", func(t *testing.T) {
			artist, err := db.GetArtists(ctx, database.ByID(artists[0].ID))
			require.Error(t, err)
			assert.True(t, errors.Is(err, database.ErrNotFound), err.Error())

			assert.Nil(t, artist)
		})
	})
}
