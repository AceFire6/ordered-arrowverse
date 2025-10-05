package httpserver

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"
)

// Config defines the server config. It is read from environment variables
type Config struct {
	Port               string        `env:"PORT" envDefault:"8080"`
	DemoModeHostPrefix string        `env:"DEMO_MODE_HOST_PREFIX,notEmpty" envDefault:"demo."`
	ShutdownTimeout    time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`
	LogLevel           zerolog.Level `env:"LOG_LEVEL" envDefault:"INFO"`
}

// GetConfig collects and returns [httpserver.Config] using environment variables or an error if values are missing
func GetConfig(devMode bool) (*Config, error) {
	cfg := &Config{}

	opts := env.Options{
		RequiredIfNoDef: true,
	}

	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return cfg, fmt.Errorf("parsing server config: %w", err)
	}

	if devMode {
		cfg.LogLevel = zerolog.DebugLevel
		cfg.ShutdownTimeout = 200 * time.Millisecond
	}

	return cfg, nil
}
