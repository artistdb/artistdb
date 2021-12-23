package database

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/johejo/golang-migrate-extra/source/iofs"
	"go.uber.org/zap"

	"github.com/obitech/artist-db/internal/database/artist"
	"github.com/obitech/artist-db/internal/database/core"
	"github.com/obitech/artist-db/internal/database/location"
)

// Database allows interaction with the underlying Postgres.
type Database struct {
	ArtistHandler   *artist.Handler
	LocationHandler *location.Handler

	conn   core.Connection
	logger *zap.Logger
}

// NewDatabase returns a database with an active connection pool.
func NewDatabase(ctx context.Context, connString string, logger *zap.Logger) (*Database, error) {
	conn, err := core.NewConnectionPool(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("connecting to DB failed: %w", err)
	}

	return &Database{
		ArtistHandler:   artist.NewHandler(conn, logger),
		LocationHandler: location.NewHandler(conn, logger),
		conn:            conn,
		logger:          logger,
	}, nil
}

// Ready returns nil if a connection to the database can be established.
func (db *Database) Ready(ctx context.Context) error {
	return db.conn.Ping(ctx)
}

func (db *Database) Close() {
	db.conn.Close()
}

//go:embed migrations/*.sql
var fs embed.FS

// TODO: migrate to in-tree iofs after
//  https://github.com/golang-migrate/migrate/issues/629 is resolved

// CreateTables creates the database tables from migration scripts.
func (db *Database) CreateTables(connString string) error {
	d, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("creating migrations dir failed: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, connString)
	if err != nil {
		return fmt.Errorf("loading migration scripts failed: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("running migrations failed: %w", err)
	}

	return nil
}

// DestroyTables runs the migration scripts backwards
func (db *Database) DestroyTables(connString string) error {
	d, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("creating migrations dir failed: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, connString)
	if err != nil {
		return fmt.Errorf("loading migration scripts failed: %w", err)
	}

	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("running migrations failed: %w", err)
	}

	return nil
}
