package handlers

import (
	"crypto"
	"encoding/base64"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

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

	inDevelopment := environment == build.Development
	assetsFs := getFileSystem(inDevelopment)

	// Static routes
	// TODO: Decide on serving the whole directory or only specific files
	// This links to the unminified files during development
	if inDevelopment {
		e.FileFS(mainCSSURL, "css/main.css", assetsFs).Name = "static:css:main"
		e.FileFS("/static/js/htmx.min.js", "js/htmx.js", assetsFs).Name = "static:js:htmx"
	} else {
		e.FileFS(mainCSSURL, "css/main.min.css", assetsFs).Name = "static:css:main"
		e.FileFS("/static/js/htmx.min.js", "js/htmx.min.js", assetsFs).Name = "static:js:htmx"
	}
	e.FileFS("/static/css/index.css", "css/index.css", assetsFs).Name = "static:css:index"
	e.FileFS("/static/js/index.js", "js/index.js", assetsFs).Name = "static:js:index"
	e.FileFS("/static/js/htmx-class-tools.js", "js/htmx-class-tools.js", assetsFs).Name = "static:js:htmx-class-tools"
	e.FileFS("/static/js/response-targets.js", "js/response-targets.js", assetsFs).Name = "static:js:htmx-response-targets"
	e.FileFS("/static/js/hyperscript.min.js", "js/hyperscript.min.js", assetsFs).Name = "static:js:hyperscript"

	e.FileFS("/favicon.png", "favicon.png", assetsFs).Name = "static:favicon"
	e.FileFS("/adsRoute.txt", "templates/ads.txt", assetsFs).Name = "static:adsRoute.txt"
	e.FileFS("/legal/privacy-policy", "templates/privacy_policy.html", assetsFs).Name = "static:privacy-policy"
	e.FileFS("/legal/cookie-policy", "templates/cookie_policy.html", assetsFs).Name = "static:cookie-policy"

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

	log.Debug().Msg("Registered routes")
	for _, route := range e.Routes() {
		log.Debug().Msgf("(Route) [%s] %s %s", route.Name, route.Method, route.Path)
	}
}
