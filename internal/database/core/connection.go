package core

import (
	"context"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/obitech/artist-db/internal/observability"
)

const (
	commandBegin = "begin"
	commandPing  = "ping"
	commandQuery = "query"
	commandExec  = "exec"
)

// Connection abstracts a pgx Database Connection.
type Connection interface {
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

// ConnectionPool wraps a pgxpool and implements the Connection interface.
type ConnectionPool struct {
	pool   *pgxpool.Pool
	tracer *trace.TracerProvider
}

type ConnectionOption func(pool *ConnectionPool) error

func NewConnectionPool(ctx context.Context, connString string) (*ConnectionPool, error) {
	conn, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	return &ConnectionPool{pool: conn}, nil
}

func (c *ConnectionPool) Ping(ctx context.Context) error {
	start := time.Now()

	defer func(s time.Time) {
		observability.Metrics.ObserveCommandDuration(commandPing, time.Since(s))
	}(start)

	return c.pool.Ping(ctx)
}

func (c *ConnectionPool) Begin(ctx context.Context) (pgx.Tx, error) {
	start := time.Now()

	defer func(s time.Time) {
		observability.Metrics.ObserveCommandDuration(commandBegin, time.Since(s))
	}(start)

	return c.pool.Begin(ctx)
}

func (c *ConnectionPool) Close() {
	c.pool.Close()
}

func (c *ConnectionPool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	start := time.Now()

	defer func(s time.Time) {
		observability.Metrics.ObserveCommandDuration(commandQuery, time.Since(s))
	}(start)

	return c.pool.Query(ctx, sql, args...)
}

func (c *ConnectionPool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	start := time.Now()

	defer func(s time.Time) {
		observability.Metrics.ObserveCommandDuration(commandQuery, time.Since(s))
	}(start)

	return c.pool.QueryRow(ctx, sql, args...)
}

func (c *ConnectionPool) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	start := time.Now()

	defer func(s time.Time) {
		observability.Metrics.ObserveCommandDuration(commandExec, time.Since(s))
	}(start)

	return c.pool.Exec(ctx, sql, args...)
}
