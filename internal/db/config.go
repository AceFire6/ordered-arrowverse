package db

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5/tracelog"
)

type Config struct {
	DatabaseURL string            `env:"DATABASE_URL,expand,notEmpty"`
	LogLevel    tracelog.LogLevel `env:"DATABASE_LOG_LEVEL" envDefault:"WARN"`
}

func GetConfig() (*Config, error) {
	cfg := &Config{
		DatabaseURL: "",
		LogLevel:    tracelog.LogLevelWarn,
	}

	opts := env.Options{
		RequiredIfNoDef: true,
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(tracelog.LogLevelInfo): func(v string) (interface{}, error) {
				return tracelog.LogLevelFromString(strings.ToLower(v))
			},
		},
	}

	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return cfg, fmt.Errorf("parsing database config: %w", err)
	}

	return cfg, nil
}
