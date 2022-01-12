package core

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/obitech/artist-db/internal/observability"
)

func RollbackAndLogError(ctx context.Context, tx pgx.Tx, logger *zap.Logger) {
	if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		observability.Metrics.TrackCommandError(commandRollback)
		logger.Error("close failed", zap.Error(err))
	}
}
