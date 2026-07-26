package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/AceFire6/ordered-arrowverse/internal/customctx"
)

// atomFeed models the wire shape of an Atom 1.0 feed document.
type atomFeed struct {
	XMLName xml.Name    `xml:"http://www.w3.org/2005/Atom feed"`
	Title   string      `xml:"title"`
	ID      string      `xml:"id"`
	Links   []atomLink  `xml:"link"`
	Logo    string      `xml:"logo"`
	Icon    string      `xml:"icon"`
	Updated string      `xml:"updated"`
	Entries []atomEntry `xml:"entry"`
}

type atomLink struct {
	XMLName xml.Name `xml:"link"`
	HREF    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr,omitempty"`
}

type atomEntry struct {
	Title   string    `xml:"title"`
	ID      string    `xml:"id"`
	Link    atomLink  `xml:"link"`
	Updated string    `xml:"updated"`
	Content atomBody  `xml:"content"`
	Author  atomActor `xml:"author"`
}

type atomBody struct {
	Type string `xml:"type,attr"`
	Text string `xml:",chardata"`
}

type atomActor struct {
	URI string `xml:"uri"`
}

// AtomFeed serves the recent episodes Atom feed. Honours the same
// `hide_show` filter as the home view, then emits the 15 newest (by
// ordering) episodes left over. Content-Type is `application/atom+xml`.
func AtomFeed(c echo.Context) error {
	cc := c.(*customctx.Context)

	var pageOpts PageOptions
	if err := c.Bind(&pageOpts); err != nil {
		cc.Log.Err(err).Msg("failed to bind atom feed page options")
		return err
	}

	cc.Log.Debug().Interface("pageOpts", pageOpts).Msg("atom feed page options")

	episodes, err := cc.ShowDB.GetEpisodesFiltered(
		c.Request().Context(),
		pageOpts.HideShowsList,
		nil,
		nil,
	)
	if err != nil {
		cc.Log.Err(err).Msg("could not load episodes for atom feed")
		return err
	}

	showRows, err := cc.ShowDB.GetShowList(c.Request().Context())
	if err != nil {
		cc.Log.Err(err).Msg("could not load show list for atom feed")
		return err
	}

	showSource := make(map[string]string, len(showRows))
	for _, sd := range showRows {
		if len(sd.DataSources) > 0 {
			showSource[sd.Slug] = sd.DataSources[0]
		}
	}

	// Newest first. The query orders ASC by `ordering` (oldest first), so
	// reverse then clip. Same as the pre-rewrite behaviour.
	slices.Reverse(episodes)
	const feedLimit = 15
	if len(episodes) > feedLimit {
		episodes = episodes[:feedLimit]
	}

	root := schemeHostRoot(c.Request())
	requestURL := root + c.Request().RequestURI

	entries := make([]atomEntry, 0, len(episodes))
	var latest time.Time
	for _, ep := range episodes {
		if ep.AirDate.After(latest) {
			latest = ep.AirDate
		}
		source := showSource[ep.ShowSlug]
		if source == "" {
			source = root
		}

		updated := ep.AirDate.UTC().Format(time.RFC3339)
		episodeID := fmt.Sprintf("S%02dE%02d", ep.Season, ep.Episode)
		title := ep.ShowName + " - " + episodeID + " - " + ep.Name
		content := ep.ShowName + " " + episodeID + " " + ep.Name +
			" will air on " + ep.AirDate.Format("January 2, 2006")

		entries = append(entries, atomEntry{
			Title:   title,
			ID:      source,
			Link:    atomLink{HREF: source},
			Updated: updated,
			Content: atomBody{Type: "text", Text: content},
			Author:  atomActor{URI: root},
		})
	}

	updated := time.Now().UTC().Format(time.RFC3339)
	if !latest.IsZero() {
		updated = latest.UTC().Format(time.RFC3339)
	}

	logoLink := root + "/favicon.png"
	feed := atomFeed{
		Title: "Arrowverse.info - Recent Episodes",
		ID:    requestURL,
		Links: []atomLink{
			{HREF: requestURL, Rel: "self"},
			{HREF: logoLink, Rel: "related"},
		},
		Logo:    logoLink,
		Icon:    logoLink,
		Updated: updated,
		Entries: entries,
	}

	body, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		cc.Log.Err(err).Msg("could not marshal atom feed")
		return err
	}
	body = append([]byte(xml.Header), body...)

	return c.Blob(http.StatusOK, "application/atom+xml; charset=utf-8", body)
}

// schemeHostRoot builds the `scheme://host` URL root for the current
// request, falling back to https when the upstream didn't send a TLS
// indicator.
func schemeHostRoot(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	return scheme + "://" + r.Host
}
