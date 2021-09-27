package server

import "go.uber.org/zap"

// Option allows customization of the default Server.
type Option func(s *Server) error

// WithLogger allows setting a custom zap logger.
func WithLogger(logger *zap.Logger) Option {
	return func(s *Server) error {
		s.logger = logger
		return nil
	}
}
