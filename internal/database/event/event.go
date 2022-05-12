package event

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID         string
	Name       string
	StartTime  *time.Time
	LocationID *string
}

type Option func(*Event)

func WithLocationID(id string) Option {
	return func(e *Event) {
		e.LocationID = &id
	}
}

func WithStartTime(startTime time.Time) Option {
	return func(e *Event) {
		t := startTime.UTC()
		e.StartTime = &t
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
