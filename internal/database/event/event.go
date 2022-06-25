package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
)

type Event struct {
	ID             string
	Name           string
	StartTime      *time.Time
	LocationID     *string
	InvitedArtists InvitedArtists
}

func (e *Event) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("id", e.ID)
	enc.AddString("name", e.Name)

	if e.StartTime != nil {
		enc.AddTime("startTime", *e.StartTime)
	}

	if e.LocationID != nil {
		enc.AddString("locationID", *e.LocationID)
	}

	return enc.AddArray("invitedArtists", e.InvitedArtists)
}

type InvitedArtists []InvitedArtist

type InvitedArtist struct {
	ID        string
	Confirmed bool
}

func (ia InvitedArtists) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, a := range ia {
		if err := enc.AppendObject(a); err != nil {
			return err
		}
	}

	return nil
}

func (i InvitedArtist) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("id", i.ID)
	enc.AddBool("confirmed", i.Confirmed)

	return nil
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
		for _, a := range artists {
			if _, err := uuid.Parse(a.ID); err != nil {
				return fmt.Errorf("invalid UUID %q: %w", a.ID, err)
			}

			e.InvitedArtists = append(e.InvitedArtists, a)
		}

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
