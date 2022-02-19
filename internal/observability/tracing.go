package observability

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"

	"github.com/go-logr/zapr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	otelTrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"

	"github.com/obitech/artist-db/internal"
	"github.com/obitech/artist-db/internal/config"
)

type tracerProvider struct {
	*trace.TracerProvider
}

type zapErrorHandler struct {
	*zap.Logger
}

func (z *zapErrorHandler) Handle(err error) {
	z.Error("otel tracing error", zap.Error(err))
}

// NewResource returns an OpenTelemetry Resource.
func newResource() (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(internal.Name),
			semconv.ServiceVersionKey.String(internal.Version),
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
func NewTracerProvider(ctx context.Context, cfg *config.Config, opts ...trace.TracerProviderOption) (*tracerProvider, error) {
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	if cfg.Tracing.SampleRate <= 0 {
		nop := otelTrace.NewNoopTracerProvider()
		otel.SetTracerProvider(nop)
		return &tracerProvider{}, nil
	}

	var (
		exp trace.SpanExporter
		err error
	)
	res, err := newResource()
	if err != nil {
		return nil, fmt.Errorf("creating resource failed: %w", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(
			trace.ParentBased(
				trace.TraceIDRatioBased(cfg.Tracing.SampleRate),
			),
		),
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)

	otel.SetErrorHandler(&zapErrorHandler{Logger: zap.L()})
	otel.SetLogger(zapr.NewLogger(zap.L()))
	otel.SetTracerProvider(tp)

	return &tracerProvider{tp}, nil
}

func ExtractTraceID(ctx context.Context) string {
	return otelTrace.SpanFromContext(ctx).SpanContext().TraceID().String()
}

func TraceField(ctx context.Context) zap.Field {
	return zap.String("trace.id", ExtractTraceID(ctx))
}
