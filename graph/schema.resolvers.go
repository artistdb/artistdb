package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/obitech/artist-db/graph/generated"
	model_gen "github.com/obitech/artist-db/graph/model"
	"github.com/obitech/artist-db/internal/database/model"
)

func (r *mutationResolver) UpsertArtists(ctx context.Context, input []*model_gen.ArtistInput) ([]*model_gen.Artist, error) {
	artists := []*model.Artist{}
	ret := []*model_gen.Artist{} //return array as workaround for now, to satisfy resolver's need to use generated model

	layout := "2020-01-01 09:09:09"

	for i := 0; i < len(input); i++ {

		var artist model.Artist
		var re model_gen.Artist

		if input[i].ID != nil {
			artist.ID = *input[i].ID
		} else {
			artist.ID = uuid.NewString()
		}
		
		artist.FirstName = input[i].FirstName
		artist.LastName = input[i].LastName
		artist.ArtistName = toString(input[i].ArtistName)
		for j := 0; j < len(input[i].Pronouns); j++ {
			artist.Pronouns[j] = toString(input[i].Pronouns[j])
		}
		artist.Origin.DateOfBirth, _ = time.Parse(layout, toString(input[i].DateOfBirth))
		artist.Origin.PlaceOfBirth = toString(input[i].PlaceOfBirth)
		artist.Language = toString(input[i].Language)
		artist.Socials.Facebook = toString(input[i].Facebook)
		artist.Socials.Instagram = toString(input[i].Instagram)
		artist.Socials.Bandcamp = toString(input[i].Bandcamp)
		artist.BioGerman = toString(input[i].BioGer)
		artist.BioEnglish = toString(input[i].BioEn)

		artists = append(artists, &artist)

		// returning a stub for now since we will refactor anyways?
		re.ID = artist.ID
		re.FirstName = artist.FirstName
		re.LastName = artist.LastName

		ret = append(ret, &re)
	}

	r.DB.UpsertArtists(ctx, artists...)

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
