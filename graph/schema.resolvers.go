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
)

func (r *mutationResolver) UpsertArtists(ctx context.Context, input []*model.ArtistInput) ([]*model.Artist, error) {
	dbArtists, err := databaseArtists(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	if err := r.db.ArtistHandler.Upsert(ctx, dbArtists...); err != nil {
		msg := "upsert failed"

		r.logger.Error(msg, zap.Error(err))
		return nil, fmt.Errorf("%s: %w", msg, err)
	}

	ret, err := modelArtists(dbArtists)
	if err != nil {
		msg := "conversion failed"

		r.logger.Error(msg, zap.Error(err))
		return nil, fmt.Errorf("%s: %w", msg, err)
	}

	return ret, nil
}

func (r *mutationResolver) DeleteArtistByID(ctx context.Context, id string) (bool, error) {
	if err := r.db.ArtistHandler.DeleteByID(ctx, id); err != nil {
		r.logger.Error("delete failed", zap.Error(err), zap.String("id", id))
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) UpsertLocations(ctx context.Context, input []*model.LocationInput) ([]*model.Location, error) {
	dbLocations, err := databaseLocations(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	if err := r.db.LocationHandler.Upsert(ctx, dbLocations...); err != nil {
		msg := "upsert failed"

		r.logger.Error(msg, zap.Error(err))
		return nil, fmt.Errorf("%s: %w", msg, err)
	}

	ret, err := modelLocations(dbLocations)
	if err != nil {
		msg := "conversion failed"

		r.logger.Error(msg, zap.Error(err))
		return nil, fmt.Errorf("%s: %w", msg, err)
	}

	return ret, nil
}

func (r *queryResolver) GetArtists(ctx context.Context, input []*model.GetArtistInput) ([]*model.Artist, error) {
	artists := make([]*artist.Artist, len(input))

	for i := range input {
		var (
			a   []*artist.Artist
			err error
		)

		switch {
		case input[i].ID != nil:
			a, err = r.db.ArtistHandler.Get(ctx, artist.ByID(*input[i].ID))
		case input[i].LastName != nil:
			a, err = r.db.ArtistHandler.Get(ctx, artist.ByLastName(*input[i].LastName))
		case input[i].ArtistName != nil:
			a, err = r.db.ArtistHandler.Get(ctx, artist.ByLastName(*input[i].ArtistName))
		}

		if err != nil {
			r.logger.Error("get failed", zap.Error(err))
			return nil, err
		}

		artists = append(artists, a...)
	}

	return modelArtists(artists)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type (
	mutationResolver struct{ *Resolver }
	queryResolver    struct{ *Resolver }
)
