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

type SiteConfig struct {
	// OldSiteHost is the legacy Heroku hostname that 301-redirects to
	// NewSiteURL with a banner on top of every page.
	OldSiteHost string `env:"OLD_SITE_HOST"`
	// NewSiteURL is the canonical absolute URL promoted on the legacy
	// site banner.
	NewSiteURL string `env:"NEW_SITE_URL,notEmpty" envDefault:"https://arrowverse.info"`
}

// GetInfo collects and returns [build.Info] using environment variables or an error if values are missing
func GetInfo() (*Info, error) {
	cfg := &Info{} //nolint:exhaustruct // these values are set by [env.ParseWithOptions]

	if err := env.Parse(cfg); err != nil {
		return cfg, fmt.Errorf("reading buildinfo: %w", err)
	}

	return cfg, nil
}

// LoadSiteConfig returns the parsed [SiteConfig]. Values have defaults
// so missing variables don't fail startup; OldSiteHost stays empty
// unless explicitly set.
func LoadSiteConfig() (*SiteConfig, error) {
	cfg := &SiteConfig{} //nolint:exhaustruct // values set by env.Parse

	if err := env.Parse(cfg); err != nil {
		return cfg, fmt.Errorf("reading site config: %w", err)
	}

	return cfg, nil
}
