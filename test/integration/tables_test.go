package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/obitech/artist-db/internal/database/core"
)

func Test_TablesExistsIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, conn, teardown := setup(t, ctx)
	defer teardown(t)

	stmt := `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema='public' and table_name=$1);`

	t.Run("artists exists", func(t *testing.T) {
		var exists bool
		require.NoError(t, conn.QueryRow(ctx, stmt, core.TableArtists).Scan(&exists))

		assert.True(t, exists)
	})

	t.Run("locations exists", func(t *testing.T) {
		var exists bool
		require.NoError(t, conn.QueryRow(ctx, stmt, core.TableLocations).Scan(&exists))

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
