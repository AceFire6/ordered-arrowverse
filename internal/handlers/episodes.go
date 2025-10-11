package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/AceFire6/ordered-arrowverse/components"
	"github.com/AceFire6/ordered-arrowverse/internal/customctx"
	"github.com/AceFire6/ordered-arrowverse/internal/frontend"
)

func Home(c echo.Context) error {
	cc := c.(*customctx.Context)

	episodes, err := cc.ShowDB.GetEpisodes(c.Request().Context())
	if err != nil {
		cc.Log.Err(err).Msg("could not load episodes")
	}

	showList := map[string]components.ShowData{}

	tableRows := make([]components.TableRow, len(episodes))
	for i, episode := range episodes {
		tableRows[i] = components.TableRow{
			ShowSlug:    episode.ShowSlug,
			RowNumber:   i + 1,
			Series:      episode.ShowName,
			EpisodeId:   fmt.Sprintf("S%02dE%02d", episode.Season, episode.Episode),
			EpisodeName: episode.Name,
			AirDate:     episode.AirDate.Format("January 2, 2006"),
			SourceLink:  episode.SourceLink,
		}
	}

	return frontend.Render(c, http.StatusOK, components.Home(cc.Frontend.DefaultPageConfig, tableRows, false, showList, []string{}, "", ""))
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
