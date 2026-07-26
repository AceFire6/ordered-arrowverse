// Package pprofsrv exposes the standard net/http/pprof endpoints on a
// separate listener when DEBUG_PPROF_ADDR is set in the environment.
//
// Pulled out of cmd/ordered-arrowverse so the main binary stays small and
// the pprof server is opt-in (default off). The internal address is
// expected to be reachable only by operators (often sidecars / port
// forwards in production).
package pprofsrv

import (
	"context"
	"errors"
	"net/http"
	"net/http/pprof"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Config turns the env vars into structured configuration. The HTTP
// listener returns nil when DEBUG_PPROF_ADDR is unset.
type Config struct {
	Addr            string        `env:"DEBUG_PPROF_ADDR"`
	ShutdownTimeout time.Duration `env:"DEBUG_PPROF_TIMEOUT" envDefault:"2s"`
}

// LoadConfig reads environment variables. Always returns a non-nil
// Config; the caller checks Addr == "" to decide whether to start.
func LoadConfig() Config {
	return Config{
		Addr:            os.Getenv("DEBUG_PPROF_ADDR"),
		ShutdownTimeout: parseDurationOr(os.Getenv("DEBUG_PPROF_TIMEOUT"), 2*time.Second),
	}
}

func parseDurationOr(in string, fallback time.Duration) time.Duration {
	if in == "" {
		return fallback
	}
	d, err := time.ParseDuration(in)
	if err != nil || d <= 0 {
		return fallback
	}

	return d
}

// Server wraps the standard library mux for net/http/pprof.
type Server struct {
	httpServer *http.Server
}

// New builds the pprof server only if cfg.Addr is non-empty. Returns nil
// otherwise so the caller can skip Start.
func New(cfg Config, log *zerolog.Logger) *Server {
	if cfg.Addr == "" {
		return nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	if log != nil {
		log.Info().Str("addr", cfg.Addr).Msg("pprof server enabled")
	}

	return &Server{httpServer: server}
}

// Addr returns the listen address in use (for logging).
func (s *Server) Addr() string {
	if s == nil || s.httpServer == nil {
		return ""
	}

	return s.httpServer.Addr
}

// Start runs the listener in a goroutine. Returns immediately; callers
// use Stop to terminate.
func (s *Server) Start(log *zerolog.Logger) {
	if s == nil {
		return
	}
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			if log != nil {
				log.Err(err).Msg("pprof server failed")
			}
		}
	}()
}

// Stop performs a graceful shutdown using the configured deadline.
func (s *Server) Stop(ctx context.Context) error {
	if s == nil || s.httpServer == nil {
		return nil
	}

	cancelCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(cancelCtx) //nolint:wrapcheck
}
