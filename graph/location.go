package graph

import (
	"github.com/google/uuid"

	"github.com/obitech/artist-db/graph/model"
	"github.com/obitech/artist-db/internal/database/location"
)

// modelLocations takes Locations returned from the database and converts them to
// Locations defined in the GraphQL model.
func modelLocations(locations ...*location.Location) ([]*model.Location, error) {
	var out []*model.Location

	for _, loc := range locations {
		out = append(out, &model.Location{
			ID:   loc.ID,
			Name: loc.Name,
		})
	}

	return out, nil
}

// databaseLocations takes InputLocations as defined in the GraphQL models and
// converts them to Locations defined in the database.
func databaseLocations(locations ...*model.LocationInput) ([]*location.Location, error) {
	var out []*location.Location

	for _, loc := range locations {
		if loc == nil {
			continue
		}

		var id string
		if loc.ID != nil {
			id = *loc.ID
		} else {
			id = uuid.NewString()
		}

		out = append(out, &location.Location{
			ID:   id,
			Name: loc.Name,
		})
	}

	return out, nil
}
