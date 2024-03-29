package graph

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/obitech/artist-db/graph/model"
	"github.com/obitech/artist-db/internal/conversion"
	"github.com/obitech/artist-db/internal/database/artist"
	"github.com/obitech/artist-db/internal/database/event"
	"github.com/obitech/artist-db/internal/database/location"
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

		if ev.StartTime != nil {
			opts = append(opts, event.WithStartTime(time.Unix(int64(*ev.StartTime), 0).UTC()))
		}

		if ev.LocationID != nil {
			opts = append(opts, event.WithLocationID(*ev.LocationID))
		}

		for _, ia := range ev.InvitedArtists {
			opts = append(opts, event.WithInvitedArtists(event.InvitedArtist{
				ID:        ia.ID,
				Confirmed: ia.Confirmed,
			}))
		}

		dbEv, err := event.New(ev.Name, opts...)
		if err != nil {
			return nil, fmt.Errorf("creating event: %w", err)
		}

		if ev.ID != nil {
			if _, err := uuid.Parse(*ev.ID); err != nil {
				return nil, fmt.Errorf("invalid event ID %q: %w", *ev.ID, err)
			}

			dbEv.ID = *ev.ID
		}

		out = append(out, dbEv)
	}

	return out, nil
}

// modelEvents takes Events returned from the database and converts them to
// Events defined in the GraphQL model.
func (r *queryResolver) modelEvents(ctx context.Context, events ...*event.Event) ([]*model.Event, error) {
	var out []*model.Event

	for _, ev := range events {
		var loc *model.Location
		if ev.LocationID != nil {
			dbLocs, err := r.db.LocationHandler.Get(ctx, location.ByID(*ev.LocationID))
			if err != nil {
				return nil, fmt.Errorf("fetching location %q: %w", *ev.LocationID, err)
			}

			l, err := modelLocations(dbLocs...)
			if err != nil {
				return nil, fmt.Errorf("convertion location: %w", err)
			}

			loc = l[0]
		}

		var artists []*model.InvitedArtist
		for _, a := range ev.InvitedArtists {
			dbArtists, err := r.db.ArtistHandler.Get(ctx, artist.ByID(a.ID))
			if err != nil {
				return nil, fmt.Errorf("fetching artist %q: %w", a.ID, err)
			}

			convertedArtists, err := modelArtists(dbArtists...)
			if err != nil {
				return nil, fmt.Errorf("converting artist: %w", err)
			}

			for _, ca := range convertedArtists {
				artists = append(artists, &model.InvitedArtist{
					Artist:    ca,
					Confirmed: a.Confirmed,
				})
			}

		}

		t := int(conversion.Time(ev.StartTime).Unix())
		out = append(out, &model.Event{
			ID:        ev.ID,
			Name:      ev.Name,
			StartTime: &t,
			Location:  loc,
			Artists:   artists,
		})
	}

	return out, nil
}
