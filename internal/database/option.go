package database

import (
	"errors"

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
