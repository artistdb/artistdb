package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/obitech/artist-db/internal/database/artist"
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

	artist1 := artist.New()
	artist2 := artist.New()

	onlyName, _ := event.New("onlyName")
	withTime, _ := event.New("withTime", event.WithStartTime(time.Time{}.UTC()))
	withLoc, _ := event.New("withLocation", event.WithLocationID(loc1.ID))
	withArtist, _ := event.New("withInvitedArtist", event.WithInvitedArtists(event.InvitedArtist{ID: artist1.ID}))
	withAll, _ := event.New("withEverything",
		event.WithStartTime(time.Time{}),
		event.WithLocationID(loc2.ID),
		event.WithInvitedArtists(
			event.InvitedArtist{ID: artist1.ID},
			event.InvitedArtist{ID: artist2.ID, Confirmed: true},
		),
	)

	events := []*event.Event{
		onlyName,
		withTime,
		withLoc,
		withAll,
		withArtist,
	}

	t.Run("inserting event without existing location throws error", func(t *testing.T) {
		require.Error(t, db.EventHandler.Upsert(ctx, events[2]))
	})

	t.Run("inserting event without existing artist throws error", func(t *testing.T) {
		require.Error(t, db.EventHandler.Upsert(ctx, events[3]))
	})

	t.Run("inserting and retrieving single event without location works", func(t *testing.T) {
		t.Run("invalid ID throws error", func(t *testing.T) {
			require.Error(t, db.EventHandler.Upsert(ctx, &event.Event{ID: "foo"}))
		})

		t.Run("insert", func(t *testing.T) {
			require.NoError(t, db.EventHandler.Upsert(ctx, events[0]))
		})

		t.Run("verify", func(t *testing.T) {
			t.Run("resources are created", func(t *testing.T) {
				events, err := db.EventHandler.Get(ctx, event.ByID(events[0].ID))
				require.NoError(t, err)

				assert.Len(t, events, 1)
				assert.Equal(t, events[0], events[0])
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

		t.Run("deleting event works", func(t *testing.T) {
			require.NoError(t, db.EventHandler.DeleteByID(ctx, events[0].ID))

			res, err := db.EventHandler.Get(ctx, event.ByID(events[0].ID))
			require.ErrorIs(t, err, core.ErrNotFound)
			require.Nil(t, res)
		})
	})

	t.Run("inserting multiple events work", func(t *testing.T) {
		t.Run("insert", func(t *testing.T) {
			require.NoError(t, db.LocationHandler.Upsert(ctx, loc1, loc2))
			require.NoError(t, db.ArtistHandler.Upsert(ctx, artist1, artist2))
			require.NoError(t, db.EventHandler.Upsert(ctx, events...))
		})

		t.Run("verify", func(t *testing.T) {
			// ByID
			ev, err := db.EventHandler.Get(ctx, event.ByID(events[0].ID))
			require.NoError(t, err)

			require.Len(t, ev, 1)
			assert.Equal(t, events[0], ev[0])

			ev, err = db.EventHandler.Get(ctx, event.ByID(events[1].ID))
			require.NoError(t, err)

			require.Len(t, ev, 1)
			assert.Equal(t, events[1], ev[0])

			ev, err = db.EventHandler.Get(ctx, event.ByID(events[2].ID))
			require.NoError(t, err)

			require.Len(t, ev, 1)
			assert.Equal(t, events[2], ev[0])

			ev, err = db.EventHandler.Get(ctx, event.ByID(events[3].ID))
			require.NoError(t, err)

			require.Len(t, ev, 1)
			assert.Equal(t, events[3], ev[0])

			// ByName
			ev, err = db.EventHandler.Get(ctx, event.ByName(events[0].Name))
			require.NoError(t, err)

			require.Len(t, ev, 1)
			assert.Equal(t, events[0], ev[0])

			ev, err = db.EventHandler.Get(ctx, event.ByName(events[1].Name))
			require.NoError(t, err)

			require.Len(t, ev, 1)
			assert.Equal(t, events[1], ev[0])

			ev, err = db.EventHandler.Get(ctx, event.ByName(events[2].Name))
			require.NoError(t, err)

			require.Len(t, ev, 1)
			assert.Equal(t, events[2], ev[0])

			ev, err = db.EventHandler.Get(ctx, event.ByName(events[3].Name))
			require.NoError(t, err)

			require.Len(t, ev, 1)
			assert.Equal(t, events[3], ev[0])
		})

		t.Run("updating existing events work", func(t *testing.T) {
			newName := "newName"
			events[0].Name = newName
			require.NoError(t, db.EventHandler.Upsert(ctx, events[0]))

			ev, err := db.EventHandler.Get(ctx, event.ByName(newName))
			require.NoError(t, err)

			require.Len(t, ev, 1)
			assert.Equal(t, events[0], ev[0])
		})
	})

	t.Run("inserting event with existing, assigned location ID works", func(t *testing.T) {
		e, err := event.New("withDuplicateLocation", event.WithLocationID(loc1.ID))
		require.NoError(t, err)

		events = append(events, e)
		require.NoError(t, db.EventHandler.Upsert(ctx, events[4]))
	})

	t.Run("deleting assigned location erases location in event", func(t *testing.T) {
		locID := events[3].LocationID
		evID := events[3].ID

		t.Run("soft delete", func(t *testing.T) {
			require.NoError(t, db.LocationHandler.DeleteByID(ctx, *locID))

			ev, err := db.EventHandler.Get(ctx, event.ByID(evID))
			require.NoError(t, err)

			require.Len(t, ev, 1)
			assert.NotNil(t, ev[0].LocationID)
			assert.Equal(t, locID, ev[0].LocationID)

			_, err = db.LocationHandler.Get(ctx, location.ByID(*locID))
			require.ErrorIs(t, err, core.ErrNotFound)
		})

		t.Run("hard delete", func(t *testing.T) {
			stmt := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, core.TableLocations)
			_, err := conn.Exec(ctx, stmt, locID)
			require.NoError(t, err)

			_, err = db.LocationHandler.Get(ctx, location.ByID(*locID))
			require.ErrorIs(t, err, core.ErrNotFound)

			ev, err := db.EventHandler.Get(ctx, event.ByID(evID))
			require.NoError(t, err)

			require.Len(t, ev, 1)
			assert.Nil(t, ev[0].LocationID)
		})
	})

	t.Run("deletiing invalid event throws error", func(t *testing.T) {
		t.Run("invalid UUID", func(t *testing.T) {
			require.ErrorIs(t, db.EventHandler.DeleteByID(ctx, "foo"), core.ErrInvalidUUID)
		})

		t.Run("unknown event", func(t *testing.T) {
			require.ErrorIs(t, db.EventHandler.DeleteByID(ctx, uuid.New().String()), core.ErrNotFound)
		})
	})
}
