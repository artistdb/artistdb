package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/obitech/artist-db/internal/database/core"
	"github.com/obitech/artist-db/internal/database/event"
	"github.com/obitech/artist-db/internal/database/location"
)

func Test_EventsIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, conn, teardown := setup(t, ctx)
	defer teardown(t)

	loc1 := location.New()
	loc2 := location.New()

	events := []*event.Event{
		event.New("onlyName"),
		event.New("withTime", event.WithStartTime(time.Time{})),
		event.New("withLocation", event.WithLocation(loc1)),
		event.New("withTimeLocation", event.WithStartTime(time.Time{}), event.WithLocation(loc2)),
	}

	t.Run("inserting and retrieving single event works", func(t *testing.T) {
		t.Run("invalid ID throws error", func(t *testing.T) {
			require.Error(t, db.EventHandler.Upsert(ctx, &event.Event{ID: "foo"}))
		})

		t.Run("insert", func(t *testing.T) {
			require.NoError(t, db.EventHandler.Upsert(ctx, events[0]))
		})

		t.Run("verify", func(t *testing.T) {
			t.Run("resources are created", func(t *testing.T) {
				stmt := `SELECT id, name, start_time, location_id FROM events WHERE id = $1`
				var res event.Event

				var locationID *string

				require.NoError(t, conn.QueryRow(ctx, stmt, events[0].ID).Scan(
					&res.ID,
					&res.Name,
					&res.StartTime,
					&locationID,
				),
				)

				res.StartTime = res.StartTime.UTC()

				assert.Equal(t, events[0], &res)
				assert.Nil(t, locationID)
			})

			t.Run("metadata is set", func(t *testing.T) {
				stmt := fmt.Sprintf(`SELECT created_at, updated_at, deleted_at FROM %s WHERE id=$1`, core.TableEvents)

				var (
					createdAt time.Time
					updatedAt time.Time
					deletedAt *time.Time
				)

				require.NoError(t, conn.QueryRow(ctx, stmt, events[0].ID).Scan(&createdAt, &updatedAt, &deletedAt))

				assert.NotZero(t, createdAt)
				assert.NotZero(t, updatedAt)
				assert.Equal(t, updatedAt, createdAt)
				assert.Nil(t, deletedAt)
			})
		})

		t.Run("cleanup", func(t *testing.T) {
			stmt := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, core.TableEvents)

			_, err := conn.Exec(ctx, stmt, events[0].ID)
			require.NoError(t, err)

			stmt = fmt.Sprintf(`SELECT id from %s WHERE id=$1`, core.TableEvents)

			var id string
			require.Error(t, conn.QueryRow(ctx, stmt, events[0].ID).Scan(&id))
			assert.Empty(t, id, "")
		})
	})

	t.Run("inserting event without location throws error", func(t *testing.T) {
		require.Error(t, db.EventHandler.Upsert(ctx, events[2]))
	})

	// TODO: verify once Get is implemented
	t.Run("inserting multiple events work", func(t *testing.T) {
		t.Run("insert", func(t *testing.T) {
			require.NoError(t, db.LocationHandler.Upsert(ctx, loc1, loc2))
			require.NoError(t, db.EventHandler.Upsert(ctx, events...))
		})

		t.Run("updating existing events work", func(t *testing.T) {
			events[0].Name = "updatedName"
			require.NoError(t, db.EventHandler.Upsert(ctx, events[0]))
		})
	})

	t.Run("inserting event with existing, assigned location ID works", func(t *testing.T) {
		events = append(events, event.New("withDuplicateLocation", event.WithLocation(loc1)))
		require.NoError(t, db.EventHandler.Upsert(ctx, events[4]))
	})

	// TODO: test once Get is implemented
	t.Run("deleting assigned location erases location in event", func(t *testing.T) {
	})
}
