package echo

import (
	"github.com/labstack/echo/v4"

	"github.com/AceFire6/ordered-arrowverse/internal/build"
	"github.com/AceFire6/ordered-arrowverse/internal/echo-goccy-json"
)

type NewEchoParams struct {
	Environment build.Environment
}

func NewEchoInstance(params NewEchoParams) *echo.Echo {
	e := echo.New()

	// Put the server in Debug mode if we're in the development environment
	e.Debug = params.Environment == build.Development
	e.JSONSerializer = echojson.JSONSerializer{}

	return e
}
