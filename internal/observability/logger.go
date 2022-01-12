package observability

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/obitech/artist-db/internal"
)

// NewLogger returns an initialized zap Logger.
func NewLogger(mode string) (*zap.Logger, error) {
	var cfg zap.Config

	switch strings.ToLower(mode) {
	case "dev":
		cfg = zap.NewDevelopmentConfig()
	case "prod":
		cfg = zap.NewProductionConfig()
	default:
		return nil, fmt.Errorf("unsupported logging mode %q", mode)
	}

	logger, err := cfg.Build(zap.AddCaller())
	if err != nil {
		return nil, fmt.Errorf("initializing logger failed: %w", err)
	}

	logger = logger.With(zap.String("service.version", internal.Version))

	zap.RedirectStdLog(logger)
	zap.ReplaceGlobals(logger)

	return logger, nil
}
