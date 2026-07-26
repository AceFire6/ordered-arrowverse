package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AceFire6/ordered-arrowverse/components"
	"github.com/AceFire6/ordered-arrowverse/internal/db"
	"github.com/AceFire6/ordered-arrowverse/internal/frontend"
	"github.com/AceFire6/ordered-arrowverse/internal/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/AceFire6/ordered-arrowverse/internal/build"
	"github.com/AceFire6/ordered-arrowverse/internal/echo"
	"github.com/AceFire6/ordered-arrowverse/internal/healthcheck"
	"github.com/AceFire6/ordered-arrowverse/internal/httpserver"
	"github.com/AceFire6/ordered-arrowverse/internal/logger"
	"github.com/AceFire6/ordered-arrowverse/internal/pprofsrv"
	"github.com/AceFire6/ordered-arrowverse/internal/redisx"
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

	// Load .env files from any of the well-known paths before reading
	// any environment-driven config. Existing real environment variables
	// are left untouched so container deployments keep working.
	if err := godotenv.Overload(".env", ".local.env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Warn().Err(err).Msg("Failed to load .env/.local.env")
	}

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

	siteConfig, err := build.LoadSiteConfig()
	if err != nil {
		log.Logger.Err(err).Msg("could not load site config")
		return 1
	}
	log.Logger.Info().
		Str("old_site_host", siteConfig.OldSiteHost).
		Str("new_site_url", siteConfig.NewSiteURL).
		Msg("Site config loaded")

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

	redisConfig, err := redisx.GetConfig()
	if err != nil {
		appLogger.Err(err).Msg("could not load redis config")
		return 1
	}
	redisClient, err := redisx.CreateClient(appLogger, redisConfig)
	if err != nil {
		appLogger.Err(err).Msg("could not connect to redis")
		return 1
	}
	defer redisx.Close(redisClient, appLogger)

	echoApp := echo.NewEchoInstance(echo.NewEchoParams{
		Environment: buildInfo.Environment,
	})

	frontendConfig, err := frontend.GetConfig()
	if err != nil {
		appLogger.Err(err).Msg("could not load frontend config")
		return 1
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	showList, err := getShowList(timeoutCtx, dbPool)
	if err != nil {
		appLogger.Err(err).Msg("could not load show list")
		return 1
	}

	// This uses the echo apps reverse function
	appFrontend := frontend.New(frontend.NewParams{
		Reverser: echoApp,
		Title:    frontendConfig.Title,
		Heading:  frontendConfig.Heading,
		ShowList: showList,
	}).WithLayout(components.PageLayout)

	appLogger.Info().Interface("showList", appFrontend.DefaultPageConfig.ShowList).Msg("frontend default page config show list")

	// We add the custom context middle here using echo.Pre which is meant to execute before the routing stack
	echo.SetEchoMiddlewareStack(echoApp, &echo.CustomContextParams{
		Log:                appLogger,
		DBPool:             dbPool,
		Frontend:           appFrontend,
		DemoModeHostPrefix: httpServerConfig.DemoModeHostPrefix,
		Environment:        buildInfo.Environment,
		Site:               siteConfig,
	})

	handlerConfig := getHandlerConfig(dbPool, redisClient)
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

	pprofServer := pprofsrv.New(pprofsrv.LoadConfig(), appLogger)
	pprofServer.Start(appLogger)

	return runServer(httpServer, httpServerConfig, appLogger, pprofServer)
}

func runServer(httpServer *http.Server, serverConfig *httpserver.Config, logger *zerolog.Logger, pprofServer *pprofsrv.Server) int {
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

	if err := pprofServer.Stop(context.Background()); err != nil {
		logger.Err(err).Msg("pprof Shutdown")
	}

	// catching ctx.Done(). timeout of 5 seconds.
	<-shutdownCtx.Done()
	logger.Info().Msg("Server shutdown 👋")

	return 0
}

func getHandlerConfig(dbPool *pgxpool.Pool, redisClient *redis.Client) *handlers.HandlerConfig {
	healthCheckService := healthcheck.NewService(
		db.NewDBCheck(dbPool),
		redisx.NewHealthCheck(redisClient),
	)
	healthCheckHandler := healthcheck.NewHealthCheckHandler(healthCheckService)

	return &handlers.HandlerConfig{
		HealthCheckHandler: healthCheckHandler,
	}
}

func getShowList(ctx context.Context, dbPool *pgxpool.Pool) ([]frontend.ShowData, error) {
	dbShowList, err := db.New(dbPool).GetShowList(ctx)
	if err != nil {
		return nil, err
	}

	showList := db.ShowListRowsToShowData(dbShowList)

	return showList, nil
}
