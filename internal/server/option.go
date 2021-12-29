package server

import (
	"errors"

	"go.opentelemetry.io/otel/trace"
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

// WithTracerProvider adds a tracer provider.
func WithTracerProvider(tp trace.TracerProvider) Option {
	return func(s *Server) error {
		if tp == nil {
			return errors.New("tracer provider is nil")
		}

		s.tracer = tp
		return nil
	}
}
