package server

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

<<<<<<< Updated upstream
	// "github.com/obitech/artist-db/graph"
=======
	"github.com/go-chi/cors"

	"github.com/obitech/artist-db/graph"
>>>>>>> Stashed changes
	"github.com/obitech/artist-db/graph/generated"
	"github.com/obitech/artist-db/internal/database"
)

// Server holds API handlers.
type Server struct {
	router chi.Router
	DB     *database.Database
	logger *zap.Logger
}

// NewServer returns a server.
func NewServer(db *database.Database, opts ...Option) (*Server, error) {
	srv := &Server{
		router: chi.NewRouter(),
		DB:     db,
		logger: zap.L().With(zap.String("component", "server")),
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
		r.Handle("/query", gqlHandler(srv.DB))
	})

	return srv, nil
}

func gqlHandler(db *database.Database) http.HandlerFunc {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: *db}}))

	return h.ServeHTTP
}

// ListenAndServe starts the HTTP server and listens for new requests.
func (s *Server) ListenAndServe(listenAddr string) error {
	return http.ListenAndServe(listenAddr, s.router)
}
