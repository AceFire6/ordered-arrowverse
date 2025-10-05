package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

	"github.com/AceFire6/ordered-arrowverse/components"
	"github.com/AceFire6/ordered-arrowverse/internal/db"
	"github.com/AceFire6/ordered-arrowverse/internal/frontend"
	"github.com/AceFire6/ordered-arrowverse/internal/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/AceFire6/ordered-arrowverse/internal/build"
	"github.com/AceFire6/ordered-arrowverse/internal/echo"
	"github.com/AceFire6/ordered-arrowverse/internal/healthcheck"
	"github.com/AceFire6/ordered-arrowverse/internal/httpserver"
	"github.com/AceFire6/ordered-arrowverse/internal/logger"
)

func main() {
	exitCode := 0
	// If we encounter a panic - recover, log, and set exitCode to 1
	defer func() {
		panicErr := recover()
		log.Error().Msgf("panic: %v", panicErr)
		exitCode = 1
	}()
	// Exits with the set exit code
	defer func() {
		os.Exit(exitCode)
	}()

	exitCode = startServer()
}

//nolint:funlen // This sets everything up - it needs to be longer than average
func startServer() int {
	buildInfo, err := build.GetInfo()
	if err != nil {
		log.Logger.Err(err).Msg("Error getting build info")
		return 1
	}

	log.Logger.Info().Msgf("Starting")

	httpServerConfig, err := httpserver.GetConfig(buildInfo.IsDevMode())
	if err != nil {
		log.Logger.Err(err).Msg("could not load http server config")
		return 1
	}

	appLogger := logger.SetupLogger(logger.Settings{
		BuildInfo:      buildInfo,
		ServiceName:    buildInfo.ServiceName,
		Environment:    buildInfo.Environment.String(),
		IsDevelopment:  buildInfo.Environment == build.Development,
		GlobalLogLevel: httpServerConfig.LogLevel,
	})

	dbConfig, err := db.GetConfig()
	if err != nil {
		appLogger.Err(err).Msg("could not load db config")
		return 1
	}

	dbPool, err := db.CreateDBPool(appLogger, dbConfig, nil)
	if err != nil {
		appLogger.Err(err).Msg("could not create db pool")
		return 1
	}
	defer db.CloseDBPool(dbPool, appLogger)

	echoApp := echo.NewEchoInstance(echo.NewEchoParams{
		Environment: buildInfo.Environment,
	})

	frontendConfig, err := frontend.GetConfig()
	if err != nil {
		appLogger.Err(err).Msg("could not load frontend config")
		return 1
	}
	// This uses the echo apps reverse function
	appFrontend := frontend.New(frontend.NewParams{
		Reverser: echoApp,
		Title:    frontendConfig.Title,
		Heading:  frontendConfig.Heading,
	}).WithLayout(components.PageLayout)

	// We add the custom context middle here using echo.Pre which is meant to execute before the routing stack
	echo.SetEchoMiddlewareStack(echoApp, &echo.CustomContextParams{
		Log:                appLogger,
		DBPool:             dbPool,
		Frontend:           appFrontend,
		DemoModeHostPrefix: httpServerConfig.DemoModeHostPrefix,
		Environment:        buildInfo.Environment,
	})

	handlerConfig := getHandlerConfig(dbPool)
	handlers.RegisterRoutes(echoApp, handlerConfig, buildInfo.Environment, appLogger)

	httpServer := httpserver.NewHTTPServer(httpserver.NewServerParams{
		Echo:   echoApp,
		Log:    appLogger,
		Config: httpServerConfig,
		// Zero values are replaced by defaults
		ReadTimout:        0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
	})

	return runServer(httpServer, httpServerConfig, appLogger)
}

func runServer(httpServer *http.Server, serverConfig *httpserver.Config, logger *zerolog.Logger) int {
	notifyCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	serverExitCodeChan := make(chan int)

	// Start server
	go func() {
		// service connections
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msgf("Failed to start listening at '%s'", httpServer.Addr)
			serverExitCodeChan <- 1
		}

		serverExitCodeChan <- 0
	}()

	logger.Info().Msgf("Server initialised: http://%s 🚀", httpServer.Addr)

	// Wait for interrupt signal to gracefully shut down the server
	select {
	case exitCode := <-serverExitCodeChan:
		logger.Debug().Msgf("Got exit code %d", exitCode)
	case <-notifyCtx.Done():
		logger.Debug().Msg("Got done from notifyCtx")
	}

	logger.Info().Msgf("Shutdown Server with Timeout %s", serverConfig.ShutdownTimeout)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), serverConfig.ShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Err(err).Msg("Server Shutdown")
		return 1
	}

	// catching ctx.Done(). timeout of 5 seconds.
	<-shutdownCtx.Done()
	logger.Info().Msg("Server shutdown 👋")

	return 0
}

func getHandlerConfig(dbPool *pgxpool.Pool) *handlers.HandlerConfig {
	healthCheckService := healthcheck.NewService(
		db.NewDBCheck(dbPool),
	)
	healthCheckHandler := healthcheck.NewHealthCheckHandler(healthCheckService)

	return &handlers.HandlerConfig{
		HealthCheckHandler: healthCheckHandler,
	}
}
