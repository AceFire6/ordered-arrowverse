package handlers

import (
	"io/fs"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/AceFire6/ordered-arrowverse/assets"
	"github.com/AceFire6/ordered-arrowverse/internal/build"
	"github.com/AceFire6/ordered-arrowverse/internal/healthcheck"
)

type HandlerConfig struct {
	HealthCheckHandler *healthcheck.Handler
}

func getFileSystem(useOS bool) fs.FS {
	if useOS {
		log.Debug().Msg("using live mode")
		return os.DirFS("assets")
	}

	log.Debug().Msg("using embed mode")
	fsys, err := fs.Sub(assets.AssetFiles, ".")
	if err != nil {
		panic(err)
	}

	return fsys
}

func RegisterRoutes(e *echo.Echo, handlerConfig *HandlerConfig, environment build.Environment, log *zerolog.Logger) {
	log.Debug().Msg("Registering routes")

	inDevelopment := environment == build.Development
	assetsFs := getFileSystem(inDevelopment)

	// Static routes
	e.FileFS("/static/css/index.css", "css/index.css", assetsFs).Name = "static:css:index"
	e.FileFS("/static/css/flatpickr.min.css", "css/flatpickr.min.css", assetsFs).Name = "static:css:flatpickr"
	e.FileFS("/static/js/htmx.min.js", "js/htmx.min.js", assetsFs).Name = "static:js:htmx"
	e.FileFS("/static/js/alpine.min.js", "js/alpine.min.js", assetsFs).Name = "static:js:alpine"
	e.FileFS("/static/js/flatpickr.min.js", "js/flatpickr.min.js", assetsFs).Name = "static:js:flatpickr"

	e.FileFS("/favicon.png", "favicon.png", assetsFs).Name = "static:favicon"
	e.FileFS("/ads.txt", "templates/ads.txt", assetsFs).Name = "static:ads.txt"
	e.GET("/legal/privacy-policy", LegalDocument("privacy")).Name = "static:privacy-policy"
	e.GET("/legal/cookie-policy", LegalDocument("cookie")).Name = "static:cookie-policy"

	// Health check routes
	healthCheck := e.Group("/health")
	healthCheck.GET("/ping", handlerConfig.HealthCheckHandler.HealthCheck).Name = "get:healthcheck:ping"

	// Home page
	e.GET("/", Home).Name = "view:home"
	e.GET("/newest_first", NewestFirst).Name = "view:newest-first"
	e.GET("/hide/:hideList", Hide).Name = "view:hide"
	e.GET("/hide/:hideList/newest_first", HideNewestFirst).Name = "view:hide:newest-first"

	// Programmatic JSON view of the filtered episode list
	e.GET("/api", API).Name = "view:api"

	// Recent episodes Atom feed (RSS subscribers)
	e.GET("/recent_episodes.atom", AtomFeed).Name = "view:rss:recent-episodes"

	log.Debug().Msg("Registered routes")
	for _, route := range e.Routes() {
		log.Debug().Msgf("(Route) [%s] %s %s", route.Name, route.Method, route.Path)
	}
}
