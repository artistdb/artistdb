package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/obitech/artist-db/internal/database"
	"github.com/obitech/artist-db/internal/server"
)

const (
	envLoggingMode = "LOGGING_MODE"
	envConnString  = "DB_CONN_STRING"
	envListenAddr  = "LISTEN"
)

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}

	return fallback
}

func main() {
	if err := initLogger(getEnv(envLoggingMode, "dev")); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger := zap.L().With(zap.String("component", "main"))

	connString := getEnv(envConnString, "")

	db, err := database.NewDatabase(ctx, connString)
	if err != nil {
		logger.Fatal("setting up database connection failed", zap.Error(err))
	}

	defer db.Close()

	if err := db.Ready(ctx); err != nil {
		logger.Fatal("database not ready", zap.Error(err))
	}

	if err := db.CreateTables(connString); err != nil {
		logger.Fatal("creating tables failed")
	}

	logger.Info("database initialized")

	srv, err := server.NewServer(db)
	if err != nil {
		logger.Fatal("setting up server failed", zap.Error(err))
	}

	listen := getEnv(envListenAddr, ":8080")
	logger.Info("Starting http server...", zap.String("listenAddress", listen))

	if err := srv.ListenAndServe(listen); err != nil {
		logger.Error("listen failed", zap.Error(err))
	}
}
