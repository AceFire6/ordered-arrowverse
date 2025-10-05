package build

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"
)

type Info struct {
	// CommitHash commit the build is based on from COMMIT_HASH
	CommitHash string `env:"COMMIT_HASH" envDefault:"DEV" json:"commit_hash"`
	// BuildDatetime datetime the build was done at from BUILD_DATETIME
	BuildDatetime time.Time `env:"BUILD_DATETIME" json:"build_datetime"`
	// ServiceName name of the service from SERVICE_NAME
	ServiceName string        `env:"SERVICE_NAME,notEmpty" json:"service_name"`
	Environment Environment   `env:"ENVIRONMENT" envDefault:"production"`
	LogLevel    zerolog.Level `env:"LOG_LEVEL" envDefault:"INFO" json:"log_level"`
}

func (i Info) IsDevMode() bool {
	return i.Environment == Development
}

// GetInfo collects and returns [build.Info] using environment variables or an error if values are missing
func GetInfo() (*Info, error) {
	cfg := &Info{} //nolint:exhaustruct // these values are set by [env.ParseWithOptions]

	if err := env.Parse(cfg); err != nil {
		return cfg, fmt.Errorf("reading buildinfo: %w", err)
	}

	return cfg, nil
}
