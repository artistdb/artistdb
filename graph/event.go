package graph

import (
	"errors"

	"github.com/obitech/artist-db/graph/model"
	"github.com/obitech/artist-db/internal/conversion"
	"github.com/obitech/artist-db/internal/database/event"
)

// databaseEvents takes EventInput as defined in the GraphQL models and
// converts them to Events defined in the database.
func databaseEvents(events ...*model.EventInput) ([]*event.Event, error) {
	var out []*event.Event
	for _, ev := range events {
		if ev == nil {
			continue
		}

		if ev.Name == "" {
			return nil, errors.New("empty name")
		}

		var opts []event.Option

		if l := ev.Location; l != nil {
			if l.ID == nil {
				return nil, errors.New("location ID is empty")
			}

			opts = append(opts, event.WithLocationID(*ev.Location.ID))
		}

		for _, ia := range ev.InvitedArtists {
			if a := ia.Artist; a == nil || a.ID == nil {
				return nil, errors.New("artist ID is empty")
			}

			opts = append(opts, event.WithInvitedArtists(event.InvitedArtist{
				ID:        *ia.Artist.ID,
				Confirmed: ia.Confirmed,
			}))
		}

		dbEv := event.New(ev.Name, opts...)
		if ev.ID != nil {
			dbEv.ID = *ev.ID
		}

		out = append(out, dbEv)
	}

	return out, nil
}

// modelEvents takes Events returned from the database and converts them to
// Events defined in the GraphQL model.
func modelEvents(events ...*event.Event) ([]*model.Event, error) {
	var out []*model.Event

	for _, ev := range events {
		var loc *model.Location
		if ev.LocationID != nil {
			// TODO: fetch location
		}

		var artists []*model.InvitedArtist
		for _, a := range ev.InvitedArtists {
			// TODO: fetch artists
			_ = a
		}

		out = append(out, &model.Event{
			ID:        ev.ID,
			Name:      ev.Name,
			StartTime: conversion.RFC3339P(conversion.Time(ev.StartTime)),
			Location:  loc,
			Artists:   artists,
		})
	}

	return out, nil
}
