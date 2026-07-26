package handlers

import (
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/AceFire6/ordered-arrowverse/components"
	"github.com/AceFire6/ordered-arrowverse/internal/customctx"
	"github.com/AceFire6/ordered-arrowverse/internal/db"
	"github.com/AceFire6/ordered-arrowverse/internal/frontend"
)

type PageOptions struct {
	FromDate      *frontend.FilterDate `query:"from_date"`
	ToDate        *frontend.FilterDate `query:"to_date"`
	HideShowsList []string             `query:"hide_show"`
	NewestFirst   bool                 `query:"newest_first"`
}

func Home(c echo.Context) error {
	cc := c.(*customctx.Context)

	var pageOpts PageOptions
	if err := c.Bind(&pageOpts); err != nil {
		cc.Log.Err(err).Msg("failed to bind page options")
		return err
	}

	cc.Log.Debug().Interface("pageOpts", pageOpts).Msg("page options")

	episodes, err := cc.ShowDB.GetEpisodesFiltered(
		c.Request().Context(),
		pageOpts.HideShowsList,
		pageOpts.FromDate.TimePtr(),
		pageOpts.ToDate.TimePtr(),
	)
	if err != nil {
		cc.Log.Err(err).Msg("could not load episodes")
		return err
	}
	tableRows := db.EpisodeRowsToTableRow(episodes)
	if pageOpts.NewestFirst {
		slices.Reverse(tableRows)
	}

	// htmx partials: when the request carries the HX-Request header
	// we only need to ship the swapped region (the episode table) so
	// the browser doesn't have to reparse the entire 500KB shell.
	// The full layout is rendered for plain browser navigations.
	if c.Request().Header.Get("HX-Request") == "true" {
		return frontend.Render(c, http.StatusOK, components.EpisodesTable(tableRows))
	}

	// use the Frontend.DefaultPageConfig.ShowList instead of loading it on each request - we can do that later if needed
	pageConfig := frontend.NewPage(components.Home(cc.Echo(), tableRows, pageOpts.NewestFirst, cc.Frontend.DefaultPageConfig.ShowList, pageOpts.HideShowsList, pageOpts.FromDate.String(), pageOpts.ToDate.String()))
	customctx.ApplySiteConfig(cc, pageConfig)
	customctx.ApplyCSRFToken(cc, pageConfig)

	return cc.Frontend.RenderPage(c, http.StatusOK, pageConfig)
}

// NewestFirst is a legacy URL that 301-redirects to the canonical Home view
// with `newest_first=true`. Replaces the pre-rewrite `/newest_first/` route.
func NewestFirst(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/?newest_first=true")
}

// Hide 301-redirects the legacy `/hide/<list>/` URL to the canonical Home
// view with each slug from the `+`-separated list attached as a
// `hide_show` query parameter.
func Hide(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, buildHomeFilterURL(c.Param("hideList"), false))
}

// HideNewestFirst 301-redirects the legacy `/hide/<list>/newest_first/`
// URL to the canonical Home view with `newest_first=true` and the
// `+`-separated hide list applied.
func HideNewestFirst(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, buildHomeFilterURL(c.Param("hideList"), true))
}

// buildHomeFilterURL constructs the canonical `/` filter URL from the
// legacy `+`-separated hide-show list and an optional newest-first flag.
// Empty hide values are dropped to avoid empty `hide_show` query params.
func buildHomeFilterURL(hideList string, newestFirst bool) string {
	query := url.Values{}

	for _, slug := range strings.Split(hideList, "+") {
		if slug == "" {
			continue
		}
		query.Add("hide_show", slug)
	}

	if newestFirst {
		query.Set("newest_first", "true")
	}

	encoded := query.Encode()
	if encoded == "" {
		return "/"
	}

	return "/?" + encoded
}

// API returns the same filtered view as Home but as a JSON document.
// Useful for downstream consumers (and the seed harvester) that need
// the filtered episode list without the surrounding HTML chrome.
func API(c echo.Context) error {
	cc := c.(*customctx.Context)

	var pageOpts PageOptions
	if err := c.Bind(&pageOpts); err != nil {
		cc.Log.Err(err).Msg("failed to bind page options")
		return err
	}

	cc.Log.Debug().Interface("pageOpts", pageOpts).Msg("api page options")

	episodes, err := cc.ShowDB.GetEpisodesFiltered(
		c.Request().Context(),
		pageOpts.HideShowsList,
		pageOpts.FromDate.TimePtr(),
		pageOpts.ToDate.TimePtr(),
	)
	if err != nil {
		cc.Log.Err(err).Msg("could not load episodes")
		return err
	}
	tableRows := db.EpisodeRowsToTableRow(episodes)
	if pageOpts.NewestFirst {
		slices.Reverse(tableRows)
	}

	return c.JSON(http.StatusOK, tableRows)
}
