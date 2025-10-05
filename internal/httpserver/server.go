package httpserver

import (
	"cmp"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type NewServerParams struct {
	Echo   *echo.Echo
	Log    *zerolog.Logger
	Config *Config
	// Optional
	ReadTimout        time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func NewHTTPServer(newSrvParams NewServerParams) *http.Server {
	srv := &http.Server{
		Addr:              net.JoinHostPort("0.0.0.0", newSrvParams.Config.Port),
		Handler:           newSrvParams.Echo,
		ReadTimeout:       cmp.Or(newSrvParams.ReadTimout, 3*time.Second),
		ReadHeaderTimeout: cmp.Or(newSrvParams.ReadHeaderTimeout, 3*time.Second),
		WriteTimeout:      cmp.Or(newSrvParams.WriteTimeout, 5*time.Second),
		IdleTimeout:       cmp.Or(newSrvParams.IdleTimeout, 10*time.Second),
	}

	return srv
}
