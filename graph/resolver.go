package graph

import (
	"go.uber.org/zap"

	"github.com/obitech/artist-db/internal/database"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	db     *database.Database
	logger *zap.Logger
}

func NewResolver(db *database.Database, logger *zap.Logger) *Resolver {
	return &Resolver{
		db:     db,
		logger: logger,
	}
}
