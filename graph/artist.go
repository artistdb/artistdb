package graph

import (
	"time"

	"github.com/google/uuid"

	"github.com/obitech/artist-db/graph/model"
	"github.com/obitech/artist-db/internal/conversion"
	"github.com/obitech/artist-db/internal/database/artist"
)

// modelArtists takes Artists returned from the database and converts them to
// Artists defined in the GraphQL model.
func modelArtists(artists ...*artist.Artist) ([]*model.Artist, error) {
	var out []*model.Artist

	for _, a := range artists {
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

// databaseArtists takes InputArtists as defined in the GraphQL models and
// converts them to Artists defined in the database.
func databaseArtists(artists ...*model.ArtistInput) ([]*artist.Artist, error) {
	var out []*artist.Artist

	for _, a := range artists {
		if a == nil {
			continue
		}

		var id string
		if a.ID != nil {
			id = *a.ID
		} else {
			id = uuid.NewString()
		}

		pronouns := make([]string, len(a.Pronouns))
		for i, pronoun := range a.Pronouns {
			pronouns[i] = conversion.String(pronoun)
		}

		var dob time.Time
		if a.DateOfBirth != nil {
			dob = time.Unix(int64(*a.DateOfBirth), 0).UTC()
		}

		out = append(out, &artist.Artist{
			ID:         id,
			FirstName:  a.FirstName,
			LastName:   a.LastName,
			ArtistName: conversion.String(a.ArtistName),
			Pronouns:   pronouns,
			Origin: artist.Origin{
				DateOfBirth:  dob,
				PlaceOfBirth: conversion.String(a.PlaceOfBirth),
				Nationality:  conversion.String(a.Language),
			},
			Language: "",
			Socials: artist.Socials{
				Instagram: conversion.String(a.Instagram),
				Facebook:  conversion.String(a.Facebook),
				Bandcamp:  conversion.String(a.Bandcamp),
			},
			BioGerman:  conversion.String(a.BioGer),
			BioEnglish: conversion.String(a.BioEn),
		})
	}

	return out, nil
}
