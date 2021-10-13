package database

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/johejo/golang-migrate-extra/source/iofs"
)

const (
	TableArtists               = "artists"
	TableLocations             = "locations"
	TableEvents                = "events"
	TableInvitedArtists        = "invited_artists"
	TableArtworks              = "artworks"
	TableArtworkEventLocations = "artwork_event_locations"
)

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
