package event

import (
	"time"

	"github.com/google/uuid"

	"github.com/obitech/artist-db/internal/database/location"
)

type Event struct {
	ID        string
	Name      string
	StartTime time.Time
	Location  *location.Location
}

type Option func(*Event)

func WithLocation(loc *location.Location) Option {
	return func(e *Event) {
		e.Location = loc
	}
}

func WithStartTime(startTime time.Time) Option {
	return func(e *Event) {
		e.StartTime = startTime
	}
}

func New(name string, options ...Option) *Event {
	e := &Event{
		ID:   uuid.New().String(),
		Name: name,
	}

	for _, option := range options {
		option(e)
	}

	return e
}
