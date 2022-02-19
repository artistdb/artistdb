package server

import (
	"errors"

	"go.uber.org/zap"
)

// Option allows customization of the default Server.
type Option func(s *Server) error

// WithLogger adds a zap logger.
func WithLogger(logger *zap.Logger) Option {
	return func(s *Server) error {
		if logger == nil {
			return errors.New("logger is nil")
		}

		s.logger = logger
		return nil
	}
}
