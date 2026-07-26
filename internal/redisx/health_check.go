package redisx

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// healthCheck is the [healthcheck.healthChecker] implementation for
// Redis. It pings the server with the same timeout the client uses
// for dials.
type healthCheck struct {
	client  *redis.Client
	timeout time.Duration
}

// NewHealthCheck returns a health-checker for the provided Redis
// client. The named output ("redis") appears in the
// `/health/ping` JSON response under that key.
func NewHealthCheck(client *redis.Client) *healthCheck {
	return &healthCheck{client: client, timeout: client.Options().DialTimeout}
}

func (hc *healthCheck) Name() string {
	return "redis"
}

func (hc *healthCheck) RunCheck(ctx context.Context) error {
	checkCtx, cancel := context.WithTimeout(ctx, hc.timeout)
	defer cancel()

	if err := hc.client.Ping(checkCtx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}

	return nil
}
