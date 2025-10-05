package frontend

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Title   string `env:"SITE_TITLE,notEmpty"`
	Heading string `env:"SITE_HEADING,notEmpty"`
}

// GetConfig collects and returns [components.Config] using environment variables or an error if values are missing
func GetConfig() (*Config, error) {
	cfg := &Config{} //nolint:exhaustruct // the values are set by the env.ParseWithOptions function

	if err := env.Parse(cfg); err != nil {
		return cfg, fmt.Errorf("parsing frontend config: %w", err)
	}

	return cfg, nil
}
