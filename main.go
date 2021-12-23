package main

import (
	"context"
	"log"
	"time"

	"go.uber.org/zap"

	"github.com/obitech/artist-db/internal/config"
	"github.com/obitech/artist-db/internal/database"
	"github.com/obitech/artist-db/internal/server"
)

func main() {
	cfg := config.New()

	logger, err := initLogger(cfg.LoggingMode)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.NewDatabase(ctx, cfg.DbConnectionString, logger)
	if err != nil {
		logger.Fatal("setting up database connection failed", zap.Error(err))
	}

	defer db.Close()
	defer logger.Sync()

	if err := db.Ready(ctx); err != nil {
		logger.Fatal("database not ready", zap.Error(err))
	}

	if err := db.CreateTables(cfg.DbConnectionString); err != nil {
		logger.Fatal("creating tables failed")
	}

	logger.Info("database initialized")

	srv, err := server.NewServer(db, server.WithLogger(logger))
	if err != nil {
		logger.Fatal("setting up server failed", zap.Error(err))
	}

	logger.Info("Starting HTTP server...", zap.String("listenAddress", cfg.ListenAddress))

	if err := srv.ListenAndServe(cfg.ListenAddress); err != nil {
		logger.Error("listen failed", zap.Error(err))
	}
}
