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
			stmt := fmt.Sprintf(`SELECT name FROM %s WHERE id=$1`, core.TableLocations)

			var name string
			require.NoError(t, conn.QueryRow(ctx, stmt, locations[0].ID).Scan(&name))

			assert.Equal(t, locations[0].Name, name)
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
}
