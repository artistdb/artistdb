package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

// Database allows interaction with the underlying Postgres.
type Database struct {
	conn   Connection
	logger *zap.Logger
}

// NewDatabase returns a database with an active connection pool.
func NewDatabase(ctx context.Context, connString string) (*Database, error) {
	conn, err := newConnectionPool(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("connecting to DB failed: %w", err)
	}

	return &Database{
		conn:   conn,
		logger: zap.L().With(zap.String("component", "database")),
	}, nil
}

// Ready returns nil if a connection to the database can be established.
func (db *Database) Ready(ctx context.Context) error {
	return db.conn.Ping(ctx)
}

func (db *Database) Close() {
	db.conn.Close()
}

func rollbackAndLogError(ctx context.Context, tx pgx.Tx) {
	if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		zap.L().Error("close failed", zap.Error(err))
	}
}
