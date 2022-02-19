package integration

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/obitech/artist-db/internal/conversion"
	"github.com/obitech/artist-db/internal/database/artist"
	"github.com/obitech/artist-db/internal/database/core"
)

func Test_ArtistsIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, conn, teardown := setup(t, ctx)
	defer teardown(t)

	artists := []*artist.Artist{
		{
			ID:         uuid.New().String(),
			FirstName:  "first",
			LastName:   "last",
			ArtistName: "artist",
			Email:      "test@foo.com",
			Pronouns:   []string{"she", "her"},
			Origin: artist.Origin{
				DateOfBirth:  time.Time{},
				PlaceOfBirth: "de",
				Nationality:  "de",
			},
			Language: "de",
			Socials: artist.Socials{
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
			Email:      "test@bar.com",
			ArtistName: "artist2",
			Pronouns:   []string{"he", "her", "him"},
			Origin: artist.Origin{
				DateOfBirth:  time.Time{}.Add(time.Hour),
				PlaceOfBirth: "en",
				Nationality:  "en",
			},
			Language: "en",
			Socials: artist.Socials{
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
			require.Error(t, db.ArtistHandler.Upsert(ctx, &artist.Artist{ID: "foo"}))
		})

		require.NoError(t, db.ArtistHandler.Upsert(ctx, artists[0]))

		t.Run("verify", func(t *testing.T) {
			t.Run("resources are created", func(t *testing.T) {
				res, err := db.ArtistHandler.Get(ctx, artist.ByID(artists[0].ID))
				require.NoError(t, err)
				require.Len(t, res, 1)
				assert.Equal(t, artists[0], res[0])

				res, err = db.ArtistHandler.Get(ctx, artist.ByArtistName(artists[0].ArtistName))
				require.NoError(t, err)
				require.Len(t, res, 1)
				assert.Equal(t, artists[0], res[0])

				res, err = db.ArtistHandler.Get(ctx, artist.ByLastName(artists[0].LastName))
				require.NoError(t, err)
				require.Len(t, res, 1)
				assert.Equal(t, artists[0], res[0])
			})

			t.Run("metadata is set", func(t *testing.T) {
				stmt := fmt.Sprintf(`SELECT created_at, updated_at, deleted_at FROM %s WHERE id=$1`, core.TableArtists)

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
			stmt := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, core.TableArtists)

			_, err := conn.Exec(ctx, stmt, artists[0].ID)
			require.NoError(t, err)

			stmt = fmt.Sprintf(`SELECT id from %s WHERE id=$1`, core.TableArtists)

			var id string
			require.Error(t, conn.QueryRow(ctx, stmt, artists[0].ID).Scan(&id))
			assert.Empty(t, id, "")
		})
	})

	t.Run("inserting multiple artists works", func(t *testing.T) {
		require.NoError(t, db.ArtistHandler.Upsert(ctx, artists...))

		t.Run("verify", func(t *testing.T) {
			res, err := db.ArtistHandler.Get(ctx, artist.ByID(artists[0].ID))
			require.NoError(t, err)
			require.Len(t, res, 1)
			assert.Equal(t, artists[0], res[0])

			res, err = db.ArtistHandler.Get(ctx, artist.ByID(artists[1].ID))
			require.NoError(t, err)
			require.Len(t, res, 1)
			assert.Equal(t, artists[1], res[0])

			res, err = db.ArtistHandler.Get(ctx, artist.ByArtistName(artists[1].ArtistName))
			require.NoError(t, err)
			require.Len(t, res, 2)
			assert.Equal(t, artists[1], res[0])
			assert.Equal(t, artists[2], res[1])

			res, err = db.ArtistHandler.Get(ctx, artist.ByLastName(artists[1].LastName))
			require.NoError(t, err)
			require.Len(t, res, 2)
			assert.Equal(t, artists[1], res[0])
			assert.Equal(t, artists[2], res[1])
		})

		t.Run("updating existing artist works", func(t *testing.T) {
			artists[0].ArtistName = "pee.age"
			require.NoError(t, db.ArtistHandler.Upsert(ctx, artists...))

			res, err := db.ArtistHandler.Get(ctx, artist.ByID(artists[0].ID))
			require.NoError(t, err)
			require.Len(t, res, 1)
			assert.Equal(t, artists[0], res[0])
			assert.Equal(t, "pee.age", res[0].ArtistName)

			t.Run("metadata is set", func(t *testing.T) {
				stmt := fmt.Sprintf(`SELECT created_at, updated_at, deleted_at FROM %s WHERE ID=$1`, core.TableArtists)

				var (
					createdAt time.Time
					updatedAt time.Time
					deletedAt *time.Time
				)

				require.NoError(t, conn.QueryRow(ctx, stmt, artists[0].ID).Scan(&createdAt, &updatedAt, &deletedAt))

				assert.NotZero(t, createdAt)
				assert.NotZero(t, updatedAt)
				assert.True(t, updatedAt.After(createdAt))
				assert.Nil(t, deletedAt)
			})
		})
	})

	t.Run("retrieving non-existent artist throws error", func(t *testing.T) {
		t.Run("invalid ID throws error", func(t *testing.T) {
			res, err := db.ArtistHandler.Get(ctx, artist.ByID("foo"))
			require.Error(t, err)
			assert.Nil(t, res)
		})

		t.Run("unknown ID throws error", func(t *testing.T) {
			artist, err := db.ArtistHandler.Get(ctx, artist.ByID(uuid.New().String()))
			require.Error(t, err)

			assert.True(t, errors.Is(err, core.ErrNotFound), err.Error())
			assert.Nil(t, artist)
		})
	})

	t.Run("deleting non-existent artist throws error", func(t *testing.T) {
		t.Run("invalid ID throws error", func(t *testing.T) {
			err := db.ArtistHandler.DeleteByID(ctx, "foo")
			require.Error(t, err)

			assert.True(t, errors.Is(err, core.ErrInvalidUUID), err.Error())
		})

		t.Run("unknown ID thwrows error", func(t *testing.T) {
			err := db.ArtistHandler.DeleteByID(ctx, uuid.New().String())
			require.Error(t, err)

			assert.True(t, errors.Is(err, core.ErrNotFound), err.Error())
		})
	})

	t.Run("deleting artist works", func(t *testing.T) {
		t.Run("delete", func(t *testing.T) {
			err := db.ArtistHandler.DeleteByID(ctx, artists[0].ID)
			require.NoError(t, err)
		})

		t.Run("validate", func(t *testing.T) {
			t.Run("artist is deleted", func(t *testing.T) {
				artist, err := db.ArtistHandler.Get(ctx, artist.ByID(artists[0].ID))
				require.Error(t, err)
				assert.True(t, errors.Is(err, core.ErrNotFound), err.Error())

				assert.Nil(t, artist)
			})

			t.Run("metadata is set", func(t *testing.T) {
				stmt := fmt.Sprintf(`SELECT created_at, updated_at, deleted_at FROM %s WHERE ID=$1`, core.TableArtists)

				var (
					createdAt time.Time
					updatedAt time.Time
					deletedAt *time.Time
				)

				require.NoError(t, conn.QueryRow(ctx, stmt, artists[0].ID).Scan(&createdAt, &updatedAt, &deletedAt))

				assert.NotZero(t, createdAt)
				assert.NotZero(t, updatedAt)
				assert.NotNil(t, deletedAt)
				assert.Equal(t, *deletedAt, updatedAt)
				assert.True(t, conversion.Time(deletedAt).After(createdAt))
			})
		})
	})
}
