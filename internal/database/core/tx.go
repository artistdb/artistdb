package core

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	otelTrace "go.opentelemetry.io/otel/trace"

	"github.com/obitech/artist-db/internal/observability"
)

type Tx struct {
	pgx.Tx
	tracer otelTrace.TracerProvider
}

func (t *Tx) Commit(ctx context.Context) error {
	start := time.Now()

	spanCtx, span := t.tracer.Tracer(TracingInstrumentationName).Start(ctx, "tx.commit")
	defer span.End()

	defer func(s time.Time) {
		span.End()
		observability.Metrics.ObserveCommandDuration(commandCommit, time.Since(s))
	}(start)

	if err := t.Tx.Commit(spanCtx); err != nil {
		span.RecordError(err)
		observability.Metrics.TrackCommandError(commandCommit)
		return err
	}

	return nil
}

func (t *Tx) Rollback(ctx context.Context) error {
	start := time.Now()

	spanCtx, span := t.tracer.Tracer(TracingInstrumentationName).Start(ctx, "tx.rollback")
	defer span.End()

	defer func(s time.Time) {
		span.End()
		observability.Metrics.ObserveCommandDuration(commandRollback, time.Since(s))
	}(start)

	if err := t.Tx.Rollback(spanCtx); err != nil {
		span.RecordError(err)
		observability.Metrics.TrackCommandError(commandRollback)
		return err
	}

	return nil
}
