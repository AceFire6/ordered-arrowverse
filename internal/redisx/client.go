// Package redisx wraps go-redis with a config type, env-driven
// initialisation, and a healthcheck that the rest of the app can plug
// into [healthcheck.Service].
//
// The package is named `redisx` (rather than `redis`) to avoid
// colliding with the upstream import alias used throughout the code.
package redisx

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// Config holds the env-driven settings required to dial Redis.
type Config struct {
	URL     string        `env:"REDIS_URL,expand" envDefault:"redis://localhost:6379"`
	Timeout time.Duration `env:"REDIS_TIMEOUT" envDefault:"3s"`
}

// GetConfig parses [Config] from environment variables. REDIS_URL is
// required; the rest have defaults.
func GetConfig() (*Config, error) {
	cfg := &Config{} //nolint:exhaustruct // values set by env.Parse

	opts := env.Options{
		RequiredIfNoDef: false,
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(time.Duration(0)): func(v string) (interface{}, error) {
				return time.ParseDuration(v)
			},
		},
	}

	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return cfg, fmt.Errorf("parsing redis config: %w", err)
	}

	return cfg, nil
}

// CreateClient dials Redis using [Config]. The client is bounded by the
// configured timeout for connect + read + write operations.
func CreateClient(log *zerolog.Logger, cfg *Config) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	opts.DialTimeout = cfg.Timeout
	opts.ReadTimeout = cfg.Timeout
	opts.WriteTimeout = cfg.Timeout

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	if log != nil {
		log.Info().Str("url", redactPassword(cfg.URL)).Msg("Redis connection established")
	}

	return client, nil
}

// Close releases the underlying client connection pool.
func Close(client *redis.Client, log *zerolog.Logger) {
	if client == nil {
		return
	}
	if log != nil {
		log.Info().Msg("Closing Redis connection")
	}
	if err := client.Close(); err != nil {
		if log != nil {
			log.Err(err).Msg("Failed to close Redis client cleanly")
		}
	}
}

// redactPassword strips the password segment out of a Redis URL before
// logging it. `redis://:secret@host:6379/0` becomes `redis://host:6379/0`.
func redactPassword(url string) string {
	at := strings.LastIndex(url, "@")
	scheme := strings.Index(url, "://")
	if at == -1 || scheme == -1 || at <= scheme+3 {
		return url
	}

	return url[:scheme+3] + url[at+1:]
}
