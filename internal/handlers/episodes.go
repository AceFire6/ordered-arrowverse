package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/AceFire6/ordered-arrowverse/components"
	"github.com/AceFire6/ordered-arrowverse/internal/customctx"
	"github.com/AceFire6/ordered-arrowverse/internal/frontend"
)

func Home(c echo.Context) error {
	cc := c.(*customctx.Context)

	dbResults, err := cc.DB.Query(c.Request().Context(), "select * from episode;")
	if err != nil {
		cc.Log.Err(err).Msg("could not load episodes")
	}

	cc.Log.Debug().Interface("dbResults", dbResults).Msg("loaded episodes")
	showList := map[string]components.ShowData{}

	return frontend.Render(c, http.StatusOK, components.Home(cc.Frontend.DefaultPageConfig, []components.TableRow{}, false, showList, []string{}, "", ""))
}

func NewestFirst(c echo.Context) error {
	return c.String(http.StatusOK, "newest")
}

func Hide(c echo.Context) error {
	return c.String(http.StatusOK, "hide")
}

func HideNewestFirst(c echo.Context) error {
	//return c.Redirect(http.StatusMovedPermanently, "hide?newest_first")
	return c.String(http.StatusOK, "hidenewest")
}
