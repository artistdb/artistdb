package core

import (
	"context"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	otelTrace "go.opentelemetry.io/otel/trace"

	"github.com/obitech/artist-db/internal/observability"
)

const (
	commandBegin = "begin"
	commandPing  = "ping"
	commandQuery = "query"
	commandExec  = "exec"
)

const (
	TracingInstrumentationName = "database"
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
	tracer otelTrace.TracerProvider
}

func NewConnectionPool(ctx context.Context, connString string, tp otelTrace.TracerProvider) (*ConnectionPool, error) {
	spanCtx, span := tp.Tracer(TracingInstrumentationName).Start(ctx, "connect")
	defer span.End()

	conn, err := pgxpool.Connect(spanCtx, connString)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &ConnectionPool{
		pool:   conn,
		tracer: tp,
	}, nil
}

func (c *ConnectionPool) Ping(ctx context.Context) error {
	start := time.Now()
	spanCtx, span := c.tracer.Tracer(TracingInstrumentationName).Start(ctx, "ping")

	defer func(s time.Time) {
		span.End()
		observability.Metrics.ObserveCommandDuration(commandPing, time.Since(s))
	}(start)

	if err := c.pool.Ping(spanCtx); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (c *ConnectionPool) Begin(ctx context.Context) (pgx.Tx, error) {
	start := time.Now()
	spanCtx, span := c.tracer.Tracer(TracingInstrumentationName).Start(ctx, "begin")

	defer func(s time.Time) {
		span.End()
		observability.Metrics.ObserveCommandDuration(commandBegin, time.Since(s))
	}(start)

	res, err := c.pool.Begin(spanCtx)
	if err != nil {
		span.RecordError(err)
	}

	return res, err
}

func (c *ConnectionPool) Close() {
	c.pool.Close()
}

func (c *ConnectionPool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	start := time.Now()
	spanCtx, span := c.tracer.Tracer(TracingInstrumentationName).Start(ctx, "query")

	defer func(s time.Time) {
		span.End()
		observability.Metrics.ObserveCommandDuration(commandQuery, time.Since(s))
	}(start)

	res, err := c.pool.Query(spanCtx, sql, args...)
	if err != nil {
		span.RecordError(err)
	}

	return res, err
}

func (c *ConnectionPool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	start := time.Now()
	spanCtx, span := c.tracer.Tracer(TracingInstrumentationName).Start(ctx, "query.row")

	defer func(s time.Time) {
		span.End()
		observability.Metrics.ObserveCommandDuration(commandQuery, time.Since(s))
	}(start)

	return c.pool.QueryRow(spanCtx, sql, args...)
}

func (c *ConnectionPool) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	start := time.Now()
	spanCtx, span := c.tracer.Tracer(TracingInstrumentationName).Start(ctx, "exec")

	defer func(s time.Time) {
		span.End()
		observability.Metrics.ObserveCommandDuration(commandExec, time.Since(s))
	}(start)

	res, err := c.pool.Exec(spanCtx, sql, args...)
	if err != nil {
		span.RecordError(err)
	}

	return res, err
}
