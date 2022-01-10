package location

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/obitech/artist-db/internal/database/core"
	"github.com/obitech/artist-db/internal/observability"
)

// GetRequest specifies the input for an  Artists query against the database.
type GetRequest func() (string, string)

// ByID requests and Artist by ID.
func ByID(id string) GetRequest {
	return func() (string, string) {
		return id, "id=$1"
	}
}

// ByName requests Locations by name.
func ByName(name string) GetRequest {
	return func() (string, string) {
		return name, "name=$1"
	}
}

func (h *Handler) Get(ctx context.Context, request GetRequest) ([]*Location, error) {
	input, whereClause := request()

	stmt := fmt.Sprintf(`
		SELECT
			id,
			name
		FROM 
			"%s"
		WHERE deleted_at IS NULL AND `, core.TableLocations,
	)

	stmt += whereClause

	rows, err := h.conn.Query(ctx, stmt, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.ErrNotFound
		}

		observability.Metrics.TrackObjectError(entityLocation, "get")
		return nil, fmt.Errorf("query failed: %w", err)
	}

	defer rows.Close()

	var locations []*Location

	for rows.Next() {
		var (
			id   string
			name string
		)

		if err := rows.Scan(
			&id,
			&name,
		); err != nil {
			observability.Metrics.TrackObjectError(entityLocation, "get")
			return nil, fmt.Errorf("scanning rows failed: %w", err)
		}

		locations = append(locations, &Location{
			ID:   id,
			Name: name,
		})
	}

	if len(locations) == 0 {
		return nil, core.ErrNotFound
	}

	observability.Metrics.TrackObjectsRetrieved(len(locations), entityLocation)

	return locations, nil
}
