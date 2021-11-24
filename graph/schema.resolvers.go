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
		artist.ArtistName = toString(artistInput.ArtistName)

		if len(artistInput.Pronouns) > 0 {
			artist.Pronouns = make([]string, len(artistInput.Pronouns))

			for i, pronoun := range artistInput.Pronouns {
				artist.Pronouns[i] = toString(pronoun)
			}
		}

		if artistInput.DateOfBirth != nil {
			artist.Origin.DateOfBirth = time.Unix(int64(*artistInput.DateOfBirth), 0).UTC()
		}

		artist.Origin.PlaceOfBirth = toString(artistInput.PlaceOfBirth)
		artist.Language = toString(artistInput.Language)
		artist.Socials.Facebook = toString(artistInput.Facebook)
		artist.Socials.Instagram = toString(artistInput.Instagram)
		artist.Socials.Bandcamp = toString(artistInput.Bandcamp)
		artist.BioGerman = toString(artistInput.BioGer)
		artist.BioEnglish = toString(artistInput.BioEn)

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

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func toString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
