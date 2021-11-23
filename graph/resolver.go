package graph

import (
	"github.com/obitech/artist-db/internal/database"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	db *database.Database
}

func NewResolver(db *database.Database) *Resolver {
	return &Resolver{db: db}
}
