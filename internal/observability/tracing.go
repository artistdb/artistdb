package observability

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	otelTrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"

	"github.com/obitech/artist-db/internal"
	"github.com/obitech/artist-db/internal/config"
)

// NewResource returns an OpenTelemetry Resource.
func newResource(name, version string) (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
			semconv.ServiceVersionKey.String(version),
		),
	)
}

// NewStdoutExporter returns a SpanExporter that exports spans to the provided
// writer.
func newStdoutExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

func newGrpcExporter(ctx context.Context, endpoint string, withInsecure bool) (trace.SpanExporter, error) {
	opts := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(endpoint)}

	if withInsecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	} else {
		opts = append(opts, otlptracegrpc.WithTLSCredentials(credentials.NewTLS(&tls.Config{})))
	}

	return otlptrace.New(ctx, otlptracegrpc.NewClient(opts...))
}

// NewTracerProvider initializes and returns a TracerProvider.
func NewTracerProvider(ctx context.Context, cfg *config.Config, opts ...trace.TracerProviderOption) (*trace.TracerProvider, error) {
	var (
		exp trace.SpanExporter
		err error
	)

	res, err := newResource(internal.Name, "dev")
	if err != nil {
		return nil, fmt.Errorf("creating resource failed: %w", err)
	}

	switch cfg.Tracing.Mode {
	case "stdout":
		exp, err = newStdoutExporter(os.Stdout)
	case "grpc":
		exp, err = newGrpcExporter(ctx, cfg.Tracing.Grpc.Endpoint, cfg.Tracing.Grpc.Insecure)
	default:
		exp, err = newStdoutExporter(io.Discard)
	}

	if err != nil {
		return nil, fmt.Errorf("creating span exporter failed: %w", err)
	}

	opts = append(opts, trace.WithBatcher(exp), trace.WithResource(res))

	tp := trace.NewTracerProvider(opts...)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{},
		),
	)

	return tp, nil
}

func ExtractTraceID(ctx context.Context) string {
	return otelTrace.SpanFromContext(ctx).SpanContext().TraceID().String()
}
