package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/obitech/artist-db/internal/database"
	"github.com/obitech/artist-db/internal/metrics"
)

// Server holds API handlers.
type Server struct {
	router  chi.Router
	db      *database.Database
	logger  *zap.Logger
	metrics *metrics.Collector
}

// NewServer returns a server.
func NewServer(db *database.Database, opts ...Option) (*Server, error) {
	srv := &Server{
		router:  chi.NewRouter(),
		db:      db,
		logger:  zap.L().With(zap.String("component", "server")),
		metrics: metrics.NewCollector(),
	}

	for _, fn := range opts {
		if err := fn(srv); err != nil {
			return nil, fmt.Errorf("applying option failed: %w", err)
		}
	}

	if err := prometheus.Register(srv.metrics); err != nil {
		return nil, fmt.Errorf("registering collector failed: %w", err)
	}

	srv.router.Route("/internal", func(r chi.Router) {
		r.Get("/health", srv.health)
		r.Handle("/metrics", promhttp.Handler())
	})

	return srv, nil
}

// ListenAndServe starts the HTTP server and listens for new requests.
func (s *Server) ListenAndServe(listenAddr string) error {
	return http.ListenAndServe(listenAddr, s.router)
}
