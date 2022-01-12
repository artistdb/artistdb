package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/obitech/artist-db/internal/database/core"
	"github.com/obitech/artist-db/internal/database/location"
)

func Test_LocationsIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, conn, teardown := setup(t, ctx)
	defer teardown(t)

	locations := []*location.Location{
		{
			ID:   uuid.New().String(),
			Name: "foo",
		},
		{
			ID:   uuid.New().String(),
			Name: "bar",
		},
	}

	t.Run("inserting single location works", func(t *testing.T) {
		t.Run("invalid ID throws error", func(t *testing.T) {
			require.Error(t, db.LocationHandler.Upsert(ctx, &location.Location{ID: "foo"}))
		})

		require.NoError(t, db.LocationHandler.Upsert(ctx, locations[0]))

		t.Run("location is created", func(t *testing.T) {
			locs, err := db.LocationHandler.Get(ctx, location.ByID(locations[0].ID))
			require.NoError(t, err)

			assert.Len(t, locs, 1)
			assert.Equal(t, locations[0].ID, locs[0].ID)
			assert.Equal(t, locations[0].Name, locs[0].Name)

			locs, err = db.LocationHandler.Get(ctx, location.ByName(locations[0].Name))
			require.NoError(t, err)

			assert.Len(t, locs, 1)
			assert.Equal(t, locations[0].ID, locs[0].ID)
			assert.Equal(t, locations[0].Name, locs[0].Name)
		})

		t.Run("metadata is set", func(t *testing.T) {
			stmt := fmt.Sprintf(`SELECT created_at, updated_at, deleted_at FROM %s WHERE id=$1`, core.TableLocations)

			var (
				createdAt time.Time
				updatedAt time.Time
				deletedAt *time.Time
			)

			require.NoError(t, conn.QueryRow(ctx, stmt, locations[0].ID).Scan(&createdAt, &updatedAt, &deletedAt))

			assert.NotZero(t, createdAt)
			assert.NotZero(t, updatedAt)
			assert.Nil(t, deletedAt)

			assert.Equal(t, updatedAt, createdAt)
		})
	})

	t.Run("deleting location works", func(t *testing.T) {
		t.Run("invalid ID throws error", func(t *testing.T) {
			require.ErrorIs(t, db.LocationHandler.DeleteByID(ctx, "foo"), core.ErrInvalidUUID)
		})

		t.Run("deleting unknown location throws error", func(t *testing.T) {
			require.ErrorIs(t, db.LocationHandler.DeleteByID(ctx, uuid.New().String()), core.ErrNotFound)
		})

		t.Run("location is created", func(t *testing.T) {
			t.Run("create", func(t *testing.T) {
				require.NoError(t, db.LocationHandler.Upsert(ctx, locations[0]))
			})

			t.Run("verify", func(t *testing.T) {
				locs, err := db.LocationHandler.Get(ctx, location.ByID(locations[0].ID))
				require.NoError(t, err)

				assert.Len(t, locs, 1)
				assert.Equal(t, locations[0].ID, locs[0].ID)
			})
		})

		t.Run("location is deleted", func(t *testing.T) {
			t.Run("delete", func(t *testing.T) {
				require.NoError(t, db.LocationHandler.DeleteByID(ctx, locations[0].ID))
			})

			t.Run("verify", func(t *testing.T) {
				locs, err := db.LocationHandler.Get(ctx, location.ByID(locations[0].ID))
				require.ErrorIs(t, err, core.ErrNotFound)

				assert.Len(t, locs, 0)
			})
		})
	})
}
