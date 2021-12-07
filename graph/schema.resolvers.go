package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/obitech/artist-db/graph/generated"
	model_gen "github.com/obitech/artist-db/graph/model"
	"github.com/obitech/artist-db/internal/conversion"
	"github.com/obitech/artist-db/internal/database"
	"github.com/obitech/artist-db/internal/database/model"
)

func (r *mutationResolver) UpsertArtists(ctx context.Context, input []*model_gen.ArtistInput) ([]*model_gen.Artist, error) {
	artists := make([]*model.Artist, len(input))
	ret := make([]*model_gen.Artist, len(input))

	for i, artistInput := range input {
		var artist model.Artist
		var re model_gen.Artist

		if artistInput.ID != nil {
			artist.ID = *artistInput.ID
		} else {
			artist.ID = uuid.NewString()
		}

		artist.FirstName = artistInput.FirstName
		artist.LastName = artistInput.LastName
		artist.ArtistName = conversion.PointerToString(artistInput.ArtistName)

		if len(artistInput.Pronouns) > 0 {
			artist.Pronouns = make([]string, len(artistInput.Pronouns))

			for i, pronoun := range artistInput.Pronouns {
				artist.Pronouns[i] = conversion.PointerToString(pronoun)
			}
		}

		if artistInput.DateOfBirth != nil {
			artist.Origin.DateOfBirth = time.Unix(int64(*artistInput.DateOfBirth), 0).UTC()
		}

		artist.Origin.PlaceOfBirth = conversion.PointerToString(artistInput.PlaceOfBirth)
		artist.Language = conversion.PointerToString(artistInput.Language)
		artist.Socials.Facebook = conversion.PointerToString(artistInput.Facebook)
		artist.Socials.Instagram = conversion.PointerToString(artistInput.Instagram)
		artist.Socials.Bandcamp = conversion.PointerToString(artistInput.Bandcamp)
		artist.BioGerman = conversion.PointerToString(artistInput.BioGer)
		artist.BioEnglish = conversion.PointerToString(artistInput.BioEn)

		artists[i] = &artist

		// returning a stub for now since we will refactor anyways?
		re.ID = artist.ID
		re.FirstName = artist.FirstName
		re.LastName = artist.LastName

		ret[i] = &re
	}

	if err := r.db.UpsertArtists(ctx, artists...); err != nil {
		return nil, fmt.Errorf("upserting artist failed: %w", err)
	}

	return ret, nil
}

func (r *mutationResolver) DeleteArtistByID(ctx context.Context, id string) (bool, error) {
	if err := r.db.DeleteArtistByID(ctx, id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) GetArtists(ctx context.Context, input []*model_gen.GetArtistInput) ([]*model_gen.Artist, error) {
	var artists []*model.Artist
	var ret []*model_gen.Artist

	for i := range input {

		var artist []*model.Artist
		var err error

		switch {
		case input[i].ID != nil:
			artist, err = r.db.GetArtists(ctx, database.ByID(*input[i].ID))
			if err != nil {
				return nil, fmt.Errorf("retrieving artist failed: %w", err)
			}
		case input[i].LastName != nil:
			artist, err = r.db.GetArtists(ctx, database.ByLastName(*input[i].LastName))
			if err != nil {
				return nil, fmt.Errorf("retrieving artist failed: %w", err)
			}
		case input[i].ArtistName != nil:
			artist, err = r.db.GetArtists(ctx, database.ByLastName(*input[i].ArtistName))
			if err != nil {
				return nil, fmt.Errorf("retrieving artist failed: %w", err)
			}
		}

		artists = append(artists, artist...)

		re, err := conversion.ArtistToGenArtist(artist)
		if err != nil {
			return nil, fmt.Errorf("conversion failed: %w", err)
		}

		ret = append(ret, re...)
	}

	return ret, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
