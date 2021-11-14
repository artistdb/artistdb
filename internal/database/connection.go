package database

import (
	"context"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/obitech/artist-db/internal/metrics"
)

const (
	commandBegin = "begin"
	commandPing  = "ping"
	commandQuery = "query"
	commandExec  = "exec"
)

// connection abstracts a pgx Database connection.
type connection interface {
	// Ping checks if the database is reachable.
	Ping(ctx context.Context) error

	// Close closes the underlying connections.
	Close()

	// Begin starts a new transaction.
	Begin(ctx context.Context) (pgx.Tx, error)

	// Query executes a query.
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)

	// QueryRow queries for a single row.
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row

	// Exec executes an SQL statement.
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

// connectionPool wraps a pgxpool and implements the connection interface.
type connectionPool struct {
	pool *pgxpool.Pool
}

// newConnectionPool returns a connectionPool.
func newConnectionPool(ctx context.Context, connString string) (connection, error) {
	conn, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	return &connectionPool{pool: conn}, nil
}

func (c *connectionPool) Ping(ctx context.Context) error {
	start := time.Now()

	defer func(s time.Time) {
		metrics.Collector.ObserveCommandDuration(commandPing, time.Since(s))
	}(start)

	return c.pool.Ping(ctx)
}

func (c *connectionPool) Begin(ctx context.Context) (pgx.Tx, error) {
	start := time.Now()

	defer func(s time.Time) {
		metrics.Collector.ObserveCommandDuration(commandBegin, time.Since(s))
	}(start)

	return c.pool.Begin(ctx)
}

func (c *connectionPool) Close() {
	c.pool.Close()
}

func (c *connectionPool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	start := time.Now()

	defer func(s time.Time) {
		metrics.Collector.ObserveCommandDuration(commandQuery, time.Since(s))
	}(start)

	return c.pool.Query(ctx, sql, args...)
}

func (c *connectionPool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	start := time.Now()

	defer func(s time.Time) {
		metrics.Collector.ObserveCommandDuration(commandQuery, time.Since(s))
	}(start)

	return c.pool.QueryRow(ctx, sql, args...)
}

func (c *connectionPool) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	start := time.Now()

	defer func(s time.Time) {
		metrics.Collector.ObserveCommandDuration(commandExec, time.Since(s))
	}(start)

	return c.pool.Exec(ctx, sql, args...)
}
