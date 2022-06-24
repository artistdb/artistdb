package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID             string
	Name           string
	StartTime      *time.Time
	LocationID     *string
	InvitedArtists []InvitedArtist
}

type InvitedArtist struct {
	ID        string
	Confirmed bool
}

type Option func(*Event) error

func WithLocationID(id string) Option {
	return func(e *Event) error {
		if _, err := uuid.Parse(id); err != nil {
			return fmt.Errorf("invalid UUID %q: %w", id, err)
		}

		e.LocationID = &id
		return nil
	}
}

func WithStartTime(startTime time.Time) Option {
	return func(e *Event) error {
		t := startTime.UTC()
		e.StartTime = &t

		return nil
	}
}

// WithInvitedArtists allows assigning artists to an event.
func WithInvitedArtists(artists ...InvitedArtist) Option {
	return func(e *Event) error {
		e.InvitedArtists = append(e.InvitedArtists, artists...)
		return nil
	}
}

func New(name string, options ...Option) (*Event, error) {
	e := &Event{
		ID:   uuid.New().String(),
		Name: name,
	}

	for _, option := range options {
		if err := option(e); err != nil {
			return nil, err
		}
	}

	return e, nil
}
