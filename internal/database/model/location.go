package model

import "github.com/google/uuid"

type Location struct {
	ID   string
	Name string
}

func NewLocation() Location {
	return Location{ID: uuid.New().String()}
}
