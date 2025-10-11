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

type Context struct {
	echo.Context

	Log      *zerolog.Logger
	DB       *pgxpool.Pool
	Frontend *frontend.Frontend
	ShowDB   db.Querier

	Environment build.Environment
	HTMX        htmx.Context
	// DemoMode indicates that we should anonymize the image providers
	DemoMode bool
}
