package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/obitech/artist-db/internal/database"
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
