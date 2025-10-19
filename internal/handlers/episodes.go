package handlers

import (
	"fmt"
	"net/http"

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

	dbShowList, err := cc.ShowDB.GetShowList(c.Request().Context())
	if err != nil {
		cc.Log.Err(err).Msg("could not load episodes")
		return err
	}
	showList := db.ShowListRowsToShowData(dbShowList)

	episodes, err := cc.ShowDB.GetEpisodes(c.Request().Context())
	if err != nil {
		cc.Log.Err(err).Msg("could not load episodes")
		return err
	}

	tableRows := make([]frontend.TableRow, len(episodes))
	for i, episode := range episodes {
		tableRows[i] = frontend.TableRow{
			ShowSlug:    episode.ShowSlug,
			RowNumber:   i + 1,
			Series:      episode.ShowName,
			EpisodeId:   fmt.Sprintf("S%02dE%02d", episode.Season, episode.Episode),
			EpisodeName: episode.Name,
			AirDate:     episode.AirDate.Format("January 2, 2006"),
			SourceLink:  episode.SourceLink,
		}
	}

	pageConfig := frontend.NewPage(components.Home(cc.Echo(), tableRows, pageOpts.NewestFirst, showList, pageOpts.HideShowsList, pageOpts.FromDate.String(), pageOpts.ToDate.String()))

	return cc.Frontend.RenderPage(c, http.StatusOK, pageConfig)
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
