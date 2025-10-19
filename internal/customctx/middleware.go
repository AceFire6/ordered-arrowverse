package customctx

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/AceFire6/ordered-arrowverse/internal/build"
	"github.com/AceFire6/ordered-arrowverse/internal/db"
	"github.com/AceFire6/ordered-arrowverse/internal/frontend"
	"github.com/AceFire6/ordered-arrowverse/internal/htmx"
)

type MiddlewareConfig struct {
	Logger             *zerolog.Logger
	DBPool             *pgxpool.Pool
	Frontend           *frontend.Frontend
	Environment        build.Environment
	DemoModeHostPrefix string
}

func Middleware(config MiddlewareConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			customContext := &Context{
				Context:     ctx,
				Environment: config.Environment,
				Log:         config.Logger,
				DB:          config.DBPool,
				ShowDB:      db.New(config.DBPool),
				Frontend:    config.Frontend,
				HTMX:        htmx.ContextFromRequestContext(ctx),
			}

			return next(customContext)
		}
	}
}
