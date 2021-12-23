package main

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// initLogger initializes a global zap logger. The logger can be accessed via
// zap.L()
func initLogger(mode string) (*zap.Logger, error) {
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

	zap.RedirectStdLog(logger)
	zap.ReplaceGlobals(logger)

	return logger, nil
}
