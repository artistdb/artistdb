package main

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// initLogger initializes a global zap logger. The logger can be accessed via
// zap.L()
func initLogger(mode string) error {
	var cfg zap.Config

	switch strings.ToLower(mode) {
	case "dev":
		cfg = zap.NewDevelopmentConfig()
	case "prod":
		cfg = zap.NewProductionConfig()
	default:
		return fmt.Errorf("unsupported logging mode %q", mode)
	}

	logger, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("initializing logger failed: %w", err)
	}

	zap.RedirectStdLog(logger)
	zap.ReplaceGlobals(logger)

	return nil
}
