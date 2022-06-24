package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/obitech/artist-db/graph/generated"
	"github.com/obitech/artist-db/graph/model"
	"github.com/obitech/artist-db/internal/database/artist"
	"github.com/obitech/artist-db/internal/database/location"
	"github.com/obitech/artist-db/internal/observability"
)

func (r *mutationResolver) UpsertArtists(ctx context.Context, input []*model.ArtistInput) ([]*model.Artist, error) {
	dbArtists, err := databaseArtists(input...)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	if err := r.db.ArtistHandler.Upsert(ctx, dbArtists...); err != nil {
		msg := "upsert failed"

		r.logger.Error(msg, zap.Error(err), observability.TraceField(ctx))
		return nil, fmt.Errorf("%s: %w", msg, err)
	}

	ret, err := modelArtists(dbArtists...)
	if err != nil {
		msg := "conversion failed"

		r.logger.Error(msg, zap.Error(err), observability.TraceField(ctx))
		return nil, fmt.Errorf("%s: %w", msg, err)
	}

	return ret, nil
}

func (r *mutationResolver) DeleteArtistByID(ctx context.Context, id string) (bool, error) {
	if err := r.db.ArtistHandler.DeleteByID(ctx, id); err != nil {
		r.logger.Error("delete failed", zap.Error(err), zap.String("id", id), observability.TraceField(ctx))
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) UpsertLocations(ctx context.Context, input []*model.LocationInput) ([]string, error) {
	dbLocations, err := databaseLocations(input...)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	if err := r.db.LocationHandler.Upsert(ctx, dbLocations...); err != nil {
		msg := "upsert failed"

		r.logger.Error(msg, zap.Error(err), observability.TraceField(ctx))
		return nil, fmt.Errorf("%s: %w", msg, err)
	}

	var ret []string
	for _, loc := range dbLocations {
		ret = append(ret, loc.ID)
	}

	return ret, nil
}

func (r *mutationResolver) DeleteLocationByID(ctx context.Context, input string) (bool, error) {
	if err := r.db.LocationHandler.DeleteByID(ctx, input); err != nil {
		r.logger.Error("delete failed", zap.Error(err), zap.String("id", input), observability.TraceField(ctx))
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) UpsertEvents(ctx context.Context, input []*model.EventInput) ([]string, error) {
	dbEvents, err := databaseEvents(input...)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	if err := r.db.EventHandler.Upsert(ctx, dbEvents...); err != nil {
		msg := "upsert failed"

		r.logger.Error(msg, zap.Error(err), observability.TraceField(ctx))
		return nil, fmt.Errorf("%s: %w", msg, err)
	}

	var ret []string
	for _, ev := range dbEvents {
		ret = append(ret, ev.ID)
	}

	return ret, nil
}

func (r *queryResolver) GetArtists(ctx context.Context, input []*model.GetArtistInput) ([]*model.Artist, error) {
	var artists []*model.Artist

	for i := range input {
		var (
			dbArtists []*artist.Artist
			err       error
		)

		switch {
		case input[i].ID != nil:
			dbArtists, err = r.db.ArtistHandler.Get(ctx, artist.ByID(*input[i].ID))
		case input[i].LastName != nil:
			dbArtists, err = r.db.ArtistHandler.Get(ctx, artist.ByLastName(*input[i].LastName))
		case input[i].ArtistName != nil:
			dbArtists, err = r.db.ArtistHandler.Get(ctx, artist.ByLastName(*input[i].ArtistName))
		}

		if err != nil {
			r.logger.Error("get failed", zap.Error(err), observability.TraceField(ctx))
			return nil, err
		}

		a, err := modelArtists(dbArtists...)
		if err != nil {
			msg := "conversion failed"
			r.logger.Error(msg, zap.Error(err), observability.TraceField(ctx))
			return nil, fmt.Errorf("%s: %w", msg, err)
		}

		artists = append(artists, a...)
	}

	return artists, nil
}

func (r *queryResolver) GetLocations(ctx context.Context, input []*model.GetLocationInput) ([]*model.Location, error) {
	var locations []*model.Location

	for i := range input {
		var (
			dbLocations []*location.Location
			err         error
		)

		switch {
		case input[i].ID != nil:
			dbLocations, err = r.db.LocationHandler.Get(ctx, location.ByID(*input[i].ID))
		case input[i].Name != nil:
			dbLocations, err = r.db.LocationHandler.Get(ctx, location.ByName(*input[i].Name))
		}

		if err != nil {
			r.logger.Error("get failed", zap.Error(err), observability.TraceField(ctx))
			return nil, err
		}

		l, err := modelLocations(dbLocations...)
		if err != nil {
			msg := "conversion failed"
			r.logger.Error(msg, zap.Error(err), observability.TraceField(ctx))
			return nil, fmt.Errorf("%s: %w", msg, err)
		}

		locations = append(locations, l...)
	}

	return locations, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type (
	mutationResolver struct{ *Resolver }
	queryResolver    struct{ *Resolver }
)
