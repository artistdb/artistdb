package location

import (
	"go.uber.org/zap"

	"github.com/obitech/artist-db/internal/database/core"
)

// Handler returns a DB Handler that operates on Locations.
type Handler struct {
	conn   core.Connection
	logger *zap.Logger
}

// NewHandler returns a handler.
func NewHandler(conn core.Connection, logger *zap.Logger) *Handler {
	return &Handler{
		conn:   conn,
		logger: logger,
	}
}
