package config

import (
	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"

	"github.com/obitech/artist-db/internal"
)

type serviceConfig struct {
	kong.Plugins
	ConfigFile kong.ConfigFlag `env:"ADB_CONFIG_FILE" help:"path to config file" default:"./configuration/local/config.yaml"`
}

type Config struct {
	ListenAddress      string        `env:"ADB_LISTEN_ADDRESS_HTTP" help:"listen address of the http server" default:":8080"`
	LoggingMode        string        `env:"ADB_LOGGING_MODE" help:"for which environment logging should be configured (dev,prod)" enum:"dev,prod" default:"dev"`
	DbConnectionString string        `env:"ADB_CONN_STRING" help:"connection string to the database"`
	Tracing            TracingConfig `embed:"" prefix:"tracing-"`
}

type TracingConfig struct {
	SampleRate float64        `env:"ADB_TRACING_SAMPLE_RATE" help:"Rate with which to sample. 0 means tracing is disabled" default:"0.0"`
	Grpc       OtlpGrpcConfig `embed:"" prefix:"grpc-"`
}

type OtlpGrpcConfig struct {
	Endpoint string `env:"ADB_TRACING_OTEL_GRPC_ENDPOINT" help:"gRPC endpoint for otel traces"`
	Insecure bool   `env:"ADB_TRACING_OTEL_GRPC_WITH_INSECURE" help:"if the connection is insecure" default:"true"`
}

func New() *Config {
	var (
		cli serviceConfig
		cfg = &Config{}
	)

	cli.Plugins = append(cli.Plugins, cfg)

	_ = kong.Parse(&cli,
		kong.Name(internal.Name),
		kong.Configuration(kongyaml.Loader),
		kong.UsageOnError(),
		kong.Vars{
			"version": internal.Version,
		},
	)

	return cfg
}
