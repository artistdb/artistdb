package location

import "github.com/google/uuid"

type Location struct {
	ID   string
	Name string
}

func New() *Location {
	return &Location{ID: uuid.New().String()}
}
