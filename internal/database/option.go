package database

import (
	"errors"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Option allows customization of the default Database.
type Option func(db *Database) error

// WithLogger adds a zap logger.
func WithLogger(logger *zap.Logger) Option {
	return func(db *Database) error {
		if logger == nil {
			return errors.New("logger is nil")
		}

		db.logger = logger
		return nil
	}
}

// WithTracerProvider adds a tracer provider.
func WithTracerProvider(tp trace.TracerProvider) Option {
	return func(db *Database) error {
		if tp == nil {
			return errors.New("tracer provider is nil")
		}

		db.tracer = tp
		return nil
	}
}
