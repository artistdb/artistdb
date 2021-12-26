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
	ListenAddress      string `env:"ADB_LISTEN_ADDRESS_HTTP" help:"listen address of the http server" default:":8080"`
	LoggingMode        string `env:"ADB_LOGGING_MODE" help:"logging mode (dev or prod)" default:"dev"`
	DbConnectionString string `env:"ADB_CONN_STRING" help:"connection string to the database"`
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
