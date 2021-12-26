package server

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/obitech/artist-db/graph"
	"github.com/obitech/artist-db/graph/generated"
	"github.com/obitech/artist-db/internal"
	"github.com/obitech/artist-db/internal/database"
)

// Server holds API handlers.
type Server struct {
	router chi.Router
	db     *database.Database
	logger *zap.Logger
}

// NewServer returns a server.
func NewServer(db *database.Database, logger *zap.Logger, opts ...Option) (*Server, error) {
	srv := &Server{
		router: chi.NewRouter(),
		db:     db,
		logger: logger,
	}

	for _, fn := range opts {
		if err := fn(srv); err != nil {
			return nil, fmt.Errorf("applying option failed: %w", err)
		}
	}

	srv.router.Route("/internal", func(r chi.Router) {
		r.Get("/health", srv.health)
		r.Get("/playground", playground.Handler("GraphQL playground", "/query"))
		r.Handle("/metrics", promhttp.Handler())
	})

	srv.router.Route("/", func(r chi.Router) {
		r.Use(
			cors.AllowAll().Handler,
			loggingMiddleware(logger),
			prometheusMiddleware,
		)

		r.Handle("/query", gqlHandler(db, logger))
	})

	return srv, nil
}

func gqlHandler(db *database.Database, logger *zap.Logger) http.HandlerFunc {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(db, logger)}))

	return h.ServeHTTP
}

func (s *Server) versionHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte(internal.Version)); err != nil {
		msg := "write failed"
		s.logger.Error(msg, zap.Error(err))

		http.Error(w, msg, http.StatusInternalServerError)
	}
}

// ListenAndServe starts the HTTP server and listens for new requests.
func (s *Server) ListenAndServe(listenAddr string) error {
	return http.ListenAndServe(listenAddr, s.router)
}
