package conversion

import (
	"time"

	"github.com/google/uuid"

	"github.com/obitech/artist-db/graph/model"
	"github.com/obitech/artist-db/internal/conversion"
	"github.com/obitech/artist-db/internal/database/artist"
)

// ModelArtists takes Artists returned from the database and converts them to
// Artists defined in the GraphQL model.
func ModelArtists(input []*artist.Artist) ([]*model.Artist, error) {
	var out []*model.Artist

	for _, a := range input {
		if a == nil {
			continue
		}

		pronouns := make([]*string, len(a.Pronouns))
		for i, pronoun := range a.Pronouns {
			pronouns[i] = &pronoun
		}

		out = append(out, &model.Artist{
			ID:           a.ID,
			FirstName:    a.FirstName,
			LastName:     a.LastName,
			ArtistName:   &a.ArtistName,
			Pronouns:     pronouns,
			DateOfBirth:  conversion.RFC3339P(a.Origin.DateOfBirth),
			PlaceOfBirth: &a.Origin.PlaceOfBirth,
			Nationality:  &a.Origin.Nationality,
			Language:     &a.Language,
			Facebook:     &a.Socials.Facebook,
			Instagram:    &a.Socials.Instagram,
			Bandcamp:     &a.Socials.Bandcamp,
			BioGer:       &a.BioGerman,
			BioEn:        &a.BioEnglish,
		})
	}

	return out, nil
}

func DatabaseArtists(artists []*model.ArtistInput) ([]*artist.Artist, error) {
	var out []*artist.Artist

	for _, artistInput := range artists {
		if artistInput == nil {
			continue
		}

		var id string
		if artistInput.ID != nil {
			id = *artistInput.ID
		} else {
			id = uuid.NewString()
		}

		pronouns := make([]string, len(artistInput.Pronouns))
		for i, pronoun := range artistInput.Pronouns {
			pronouns[i] = conversion.String(pronoun)
		}

		var dob time.Time
		if artistInput.DateOfBirth != nil {
			dob = time.Unix(int64(*artistInput.DateOfBirth), 0).UTC()
		}

		out = append(out, &artist.Artist{
			ID:         id,
			FirstName:  artistInput.FirstName,
			LastName:   artistInput.LastName,
			ArtistName: conversion.String(artistInput.ArtistName),
			Pronouns:   pronouns,
			Origin: artist.Origin{
				DateOfBirth:  dob,
				PlaceOfBirth: conversion.String(artistInput.PlaceOfBirth),
				Nationality:  conversion.String(artistInput.Language),
			},
			Language: "",
			Socials: artist.Socials{
				Instagram: conversion.String(artistInput.Instagram),
				Facebook:  conversion.String(artistInput.Facebook),
				Bandcamp:  conversion.String(artistInput.Bandcamp),
			},
			BioGerman:  conversion.String(artistInput.BioGer),
			BioEnglish: conversion.String(artistInput.BioEn),
		})
	}

	return out, nil
}
