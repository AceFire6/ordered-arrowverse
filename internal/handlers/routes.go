package handlers

import (
	"crypto"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/AceFire6/ordered-arrowverse/internal/build"
	"github.com/AceFire6/ordered-arrowverse/internal/healthcheck"
)

type HandlerConfig struct {
	HealthCheckHandler *healthcheck.Handler
}

func getMd5Base64(filePath string) string {
	cleanPath := filepath.Clean(filePath)
	fileContents, err := os.ReadFile(cleanPath)
	if err != nil {
		log.Err(err).Str("file_path", cleanPath).Msg("failed to read file")

		return ""
	}

	md5Hash := crypto.MD5.New()
	md5Hash.Write(fileContents)
	hashBytes := md5Hash.Sum(nil)

	return base64.RawURLEncoding.EncodeToString(hashBytes)
}

func RegisterRoutes(e *echo.Echo, handlerConfig *HandlerConfig, environment build.Environment, log *zerolog.Logger) {
	log.Debug().Msg("Registering routes")

	cacheBustStr := getMd5Base64("./assets/css/main.min.css")
	mainCSSURL := fmt.Sprintf("/static/css/main-%s.min.css", cacheBustStr)

	// Static routes
	// TODO: Decide on serving the whole directory or only specific files
	// This links to the unminified files during development
	if environment == build.Development {
		e.File(mainCSSURL, "assets/css/main.css").Name = "static:css:main"
		e.File("/static/js/htmx.min.js", "assets/js/htmx.js").Name = "static:js:htmx"
	} else {
		e.File(mainCSSURL, "assets/css/main.min.css").Name = "static:css:main"
		e.File("/static/js/htmx.min.js", "assets/js/htmx.min.js").Name = "static:js:htmx"
	}
	e.File("/static/css/index.css", "assets/css/index.css").Name = "static:css:index"
	e.File("/static/js/index.js", "assets/js/index.js").Name = "static:js:index"
	e.File("/static/js/htmx-class-tools.js", "assets/js/htmx-class-tools.js").Name = "static:js:htmx-class-tools"
	e.File("/static/js/response-targets.js", "assets/js/response-targets.js").Name = "static:js:htmx-response-targets"
	e.File("/static/js/hyperscript.min.js", "assets/js/hyperscript.min.js").Name = "static:js:hyperscript"

	e.File("/favicon.png", "assets/favicon.png").Name = "static:favicon"
	e.File("/adsRoute.txt", "assets/templates/ads.txt").Name = "static:adsRoute.txt"
	e.File("/legal/privacy-policy", "assets/templates/privacy_policy.html").Name = "static:privacy-policy"
	e.File("/legal/cookie-policy", "assets/templates/cookie_policy.html").Name = "static:cookie-policy"

	// Health check routes
	healthCheck := e.Group("/health")
	healthCheck.GET("/ping", handlerConfig.HealthCheckHandler.HealthCheck).Name = "get:healthcheck:ping"

	// Home page
	e.GET("/", Home).Name = "view:home"
	e.GET("/newest_first", NewestFirst).Name = "view:newest-first"
	e.GET("/hide/:hideList", Hide).Name = "view:hide"
	e.GET("/hide/:hideList/newest_first", HideNewestFirst).Name = "view:hide:newest-first"

	log.Debug().Msg("Registered routes")
	for _, route := range e.Routes() {
		log.Debug().Msgf("(Route) [%s] %s %s", route.Name, route.Method, route.Path)
	}
}
