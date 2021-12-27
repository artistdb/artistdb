package server

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/obitech/artist-db/internal/observability"
)

func loggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			start := time.Now()
			w := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)

			defer func() {
				logger.Info("request served",
					zap.String("trace.id", observability.ExtractTraceID(r.Context())),
					zap.String("network.protocol", r.Proto),
					zap.String("http.request.method", r.Method),
					zap.String("url.path", r.URL.Path),
					zap.String("url.query", r.URL.RawQuery),
					zap.String("client.ip", r.RemoteAddr),
					zap.String("user_agent.original", r.Header.Get("User-Agent")),
					zap.Int("http.response.status_code", w.Status()),
					zap.Int64("http.response.time.ms`", time.Since(start).Milliseconds()),
					zap.Int64("http.response.body.bytes", int64(w.BytesWritten())),
				)
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		w := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)
		start := time.Now()

		defer func() {
			reqCtx := chi.RouteContext(r.Context())
			// routes added via Mount get a "/*" suffix, which we have to remove
			routeLabel := strings.ReplaceAll(strings.Join(reqCtx.RoutePatterns, ""), "/*", "")

			observability.Metrics.ObserveRequestDuration(r.Method, routeLabel, strconv.Itoa(w.Status()), time.Since(start))
			observability.Metrics.ObserveRequestSize(r.Method, routeLabel, strconv.Itoa(w.Status()), float64(r.ContentLength))
			observability.Metrics.ObserveResponseSize(r.Method, routeLabel, strconv.Itoa(w.Status()), float64(w.BytesWritten()))
		}()

		next.ServeHTTP(w, r)
	})
}
