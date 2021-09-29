package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

// Database allows interaction with the underlying Postgres.
type Database struct {
	conn   *pgxpool.Pool
	logger *zap.Logger
}

// NewDatabase returns a database with an active connection pool.
func NewDatabase(ctx context.Context, connString string) (*Database, error) {
	conn, err := pgxpool.Connect(ctx, connString)
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
