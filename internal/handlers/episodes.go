package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"

	"github.com/AceFire6/ordered-arrowverse/components"
	"github.com/AceFire6/ordered-arrowverse/internal/customctx"
	"github.com/AceFire6/ordered-arrowverse/internal/db"
	"github.com/AceFire6/ordered-arrowverse/internal/frontend"
)

type FilterDate time.Time

func (fd *FilterDate) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}

	t, err := time.Parse("2006-01-02", string(text))
	if err != nil {
		return err
	}

	*fd = FilterDate(t)
	return nil
}

func (fd FilterDate) String() string {
	// Return an empty string when the filter date is an empty value
	if fd == FilterDate(time.Time{}) {
		return ""
	}

	return time.Time(fd).Format("2006-01-02")
}

func (fd FilterDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(fd.String())
}

type PageOptions struct {
	FromDate      *FilterDate `query:"from_date"`
	ToDate        *FilterDate `query:"to_date"`
	HideShowsList []string    `query:"hide_show"`
	NewestFirst   bool        `query:"newest_first"`
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

	pageConfig := frontend.NewPage(components.Home(cc.Echo(), tableRows, pageOpts.NewestFirst, showList, []string{}, pageOpts.FromDate.String(), pageOpts.ToDate.String()))

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
