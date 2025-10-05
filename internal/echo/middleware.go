package echo

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"

	"github.com/AceFire6/ordered-arrowverse/internal/build"
	"github.com/AceFire6/ordered-arrowverse/internal/customctx"
	"github.com/AceFire6/ordered-arrowverse/internal/frontend"
)

type AddContextToLogFunc func(ctx zerolog.Context) zerolog.Context

func addContextToLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		rid := ctx.Response().Header().Get(echo.HeaderXRequestID)

		cc := ctx.(*customctx.Context)
		cc.Log.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("request_id", rid).
				Str("uri", ctx.Path()).
				Str("host", ctx.Request().Host).
				Bool("demo_mode", cc.DemoMode).
				Interface("htmx", cc.HTMX)
		})

		return next(ctx)
	}
}

type CustomContextParams struct {
	Log                *zerolog.Logger
	DBPool             *pgxpool.Pool
	Frontend           *frontend.Frontend
	Environment        build.Environment
	DemoModeHostPrefix string
}

func SetEchoMiddlewareStack(e *echo.Echo, params *CustomContextParams) {
	e.Use(
		customctx.Middleware(customctx.MiddlewareConfig{
			Logger:             params.Log,
			DBPool:             params.DBPool,
			Frontend:           params.Frontend,
			Environment:        params.Environment,
			DemoModeHostPrefix: params.DemoModeHostPrefix,
		}),
		middleware.RequestID(),
		addContextToLog,
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{ //nolint:exhaustruct // no need to specify all the fields here
			LogStatus:    true,
			LogURI:       true,
			LogMethod:    true,
			LogRequestID: true, // Set request ID
			HandleError:  true,
			LogValuesFunc: func(ctx echo.Context, v middleware.RequestLoggerValues) error {
				cc := ctx.(*customctx.Context)

				// Log incoming requests. The .Err function will make this an INFO log if there is no error otherwise an Error log
				cc.Log.Err(v.Error).Str("method", v.Method).Str("uri", v.URI).Int("status", v.Status).Str("request_id", v.RequestID).Send()

				return nil
			},
		}),
		// Add default security
		middleware.Secure(),
		// Compress responses
		middleware.Gzip(),
		// Recover from panics
		middleware.Recover(),
	)
}
