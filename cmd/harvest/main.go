// Command harvest reads the live arrowverse.info /api JSON endpoint,
// groups rows by show, and emits an idempotent SQL seed file that
// merges into the `show` and `episode` tables.
//
// Usage:
//
//	harvest --output seed.sql [--api-url https://arrowverse.info/api]
//
// The default output path is `migrations/002_seed_arrowverse.sql`.
// Re-running with the same episodes is safe; rows are upserted.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

type episode struct {
	Series      string `json:"series"`
	EpisodeID   string `json:"episode_id"`
	EpisodeName string `json:"episode_name"`
	AirDate     string `json:"air_date"`
	RowNumber   int    `json:"row_number"`
}

type seasonEpisode struct {
	Season  string
	Number  string
	series  string
	episode episode
}

// Source URL templates for each show, used to populate data_sources[].
var showSources = map[string]struct {
	slug string
	root string
}{
	"Arrow":                     {slug: "arrow", root: "https://arrow.fandom.com/wiki/"},
	"Batwoman":                  {slug: "batwoman", root: "https://arrow.fandom.com/wiki/"},
	"Black Lightning":           {slug: "black-lightning", root: "https://en.wikipedia.org/wiki/"},
	"Constantine":               {slug: "constantine", root: "https://arrow.fandom.com/wiki/"},
	"The Flash":                 {slug: "flash", root: "https://arrow.fandom.com/wiki/"},
	"Freedom Fighters: The Ray": {slug: "freedom-fighters", root: "https://arrow.fandom.com/wiki/"},
	"DC's Legends of Tomorrow":  {slug: "legends", root: "https://arrow.fandom.com/wiki/"},
	"Stargirl":                  {slug: "stargirl", root: "https://en.wikipedia.org/wiki/"},
	"Supergirl":                 {slug: "supergirl", root: "https://arrow.fandom.com/wiki/"},
	"Superman & Lois":           {slug: "superman-and-lois", root: "https://arrow.fandom.com/wiki/"},
	"Vixen":                     {slug: "vixen", root: "https://arrow.fandom.com/wiki/"},
}

const (
	defaultAPIURL     = "https://arrowverse.info/api"
	defaultOutputPath = "migrations/002_seed_arrowverse.sql"
	userAgent         = "ordered-arrowverse-harvest/1.0 (https://arrowverse.info)"
)

func main() {
	apiURL := flag.String("api-url", envOr("ARROWVERSE_API_URL", defaultAPIURL), "source /api endpoint URL")
	outputPath := flag.String("output", defaultOutputPath, "path to write the generated seed SQL")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	body, err := fetchEpisodes(logger, *apiURL)
	if err != nil {
		logger.Error("fetch failed", "err", err)
		os.Exit(1)
	}

	rows, err := parseEpisodes(body)
	if err != nil {
		logger.Error("parse failed", "err", err)
		os.Exit(1)
	}

	if err := writeSeed(logger, *outputPath, rows); err != nil {
		logger.Error("write failed", "err", err)
		os.Exit(1)
	}
}

// envOr returns the value of env or fallback if unset.
func envOr(env, fallback string) string {
	if v, ok := os.LookupEnv(env); ok && v != "" {
		return v
	}
	return fallback
}

// fetchEpisodes pulls the JSON document from apiURL with the project's
// user agent (and a generous timeout). The body is returned verbatim.
func fetchEpisodes(logger *slog.Logger, apiURL string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			logger.Warn("close body", "err", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return body, nil
}

// parseEpisodes decodes the JSON array. Schema is the legacy JSON shape
// from ordering/views.py (series / episode_id / episode_name / air_date
// / row_number).
func parseEpisodes(body []byte) ([]episode, error) {
	var rows []episode
	if err := json.Unmarshal(body, &rows); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return rows, nil
}

// writeSeed groups the parsed rows by show, emits UPSERT statements for
// each show and episode, and writes a SQL file containing all of them in
// a deterministic order. The output is idempotent: re-running merges
// duplicate source_link + episode_id rows.
func writeSeed(logger *slog.Logger, outputPath string, rows []episode) error {
	shows := make(map[string]showAggregate)
	for _, ep := range rows {
		agg, ok := shows[ep.Series]
		if !ok {
			source, found := showSources[ep.Series]
			if !found {
				logger.Warn("unknown series - emitting placeholder slug", "series", ep.Series)
				fallback := sourceSlug(ep.Series)
				agg.slug = fallback.showSlug
				agg.fanRoot = fallback.root
			} else {
				agg.slug = source.slug
				agg.fanRoot = source.root
			}
		}

		season, number, _ := strings.Cut(ep.EpisodeID, "E")
		season = strings.TrimPrefix(season, "S")
		agg.episodes = append(agg.episodes, episodeWithMeta{
			raw:    ep,
			season: season,
			number: number,
		})
		shows[ep.Series] = agg
	}

	if err := os.MkdirAll(parentDir(outputPath), 0o755); err != nil {
		return err
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer func() { _ = f.Close() }()

	if err := emitHeader(f); err != nil {
		return err
	}

	seriesNames := make([]string, 0, len(shows))
	for name := range shows {
		seriesNames = append(seriesNames, name)
	}
	sort.Strings(seriesNames)

	for _, name := range seriesNames {
		agg := shows[name]
		if err := emitShow(f, name, agg); err != nil {
			return err
		}
	}
	if err := emitEpisodeInserts(f, shows, seriesNames); err != nil {
		return err
	}
	if err := emitFooter(f); err != nil {
		return err
	}

	logger.Info("wrote seed", "path", outputPath, "shows", len(seriesNames), "episodes", len(rows))

	return nil
}

type showAggregate struct {
	slug     string
	fanRoot  string
	episodes []episodeWithMeta
}

type episodeWithMeta struct {
	raw    episode
	season string
	number string
}

func emitHeader(w io.Writer) error {
	_, err := fmt.Fprintln(w, "-- migrations/002_seed_arrowverse.sql")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, "--")
	_, err = fmt.Fprintln(w, "-- Generated by `cmd/harvest`. Do not edit by hand.")

	return err
}

func emitShow(w io.Writer, name string, agg showAggregate) error {
	firstAir := agg.episodes[0].raw.AirDate
	lastAir := agg.episodes[len(agg.episodes)-1].raw.AirDate
	firstDate := truncateDate(firstAir)
	lastDate := truncateDate(lastAir)
	root := agg.fanRoot
	// data_sources carries the canonical URL(s) we want surfaced.
	dataSources := fmt.Sprintf("ARRAY['%s']::text[]", root+showSlugRelativePath(name, agg.slug))

	_, err := fmt.Fprintf(w, `
insert into show (name, slug, first_aired, last_aired, data_sources, primary_colour, secondary_colour, accent_colour)
values (%s, %s, %s::timestamptz, %s::timestamptz, %s, %s, %s, %s)
on conflict (slug) do update set
    name = excluded.name,
    last_aired = excluded.last_aired,
    data_sources = excluded.data_sources;
`,
		sqlString(name), sqlString(agg.slug), sqlString(firstDate), sqlString(lastDate), dataSources,
		sqlString("#888888"), sqlString("#444444"), sqlString("#cccccc"))
	return err
}

func emitEpisodeInserts(w io.Writer, shows map[string]showAggregate, seriesNames []string) error {
	for _, name := range seriesNames {
		agg := shows[name]
		for i, e := range agg.episodes {
			airDate := truncateDate(e.raw.AirDate)
			sourceLink := agg.fanRoot + showSlugRelativePath(name, agg.slug)
			_, err := fmt.Fprintf(w, `
insert into episode (name, air_date, season, episode, runtime_minutes, description, source_link, ordering, show_id)
select %s, %s::timestamptz, %s::smallint, %s::smallint, %s::smallint, %s, %s, %d, show.show_id
from show where show.slug = %s
on conflict (show_id, season, episode) do update set
    name = excluded.name,
    air_date = excluded.air_date,
    source_link = excluded.source_link,
    ordering = excluded.ordering;
`,
				sqlString(e.raw.EpisodeName), sqlString(airDate),
				sqlString(e.season), sqlString(e.number),
				sqlString("45"), sqlString(""), sqlString(sourceLink),
				e.raw.RowNumber, sqlString(agg.slug))
			if err != nil {
				return err
			}
			_ = i // ordering taken from row_number
		}
	}

	return nil
}

func emitFooter(w io.Writer) error {
	_, err := fmt.Fprint(w, `
-- Down migration (manual): truncate the seeded data.
-- truncate table episode restart identity cascade;
-- truncate table show restart identity cascade;
`)
	return err
}

// truncateDate strips the trailing time-of-day from an ISO-8601 string,
// e.g. "2012-10-10T00:00:00" -> "2012-10-10".
func truncateDate(in string) string {
	if idx := strings.Index(in, "T"); idx > 0 {
		return in[:idx]
	}

	return in
}

// sqlString wraps a string literal for PostgreSQL by escaping single
// quotes. Not a general SQL escaper; only used for trusted seed values.
func sqlString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

// parentDir returns the directory portion of a path, or "." if empty.
func parentDir(p string) string {
	if idx := strings.LastIndex(p, "/"); idx >= 0 {
		return p[:idx]
	}
	return "."
}

// showSlugRelativePath builds the canonical fandom URL fragment
// appropriate for each show.
func showSlugRelativePath(name, slug string) string {
	switch name {
	case "Arrow":
		return "List_of_Arrow_episodes"
	case "Batwoman":
		return "List_of_Batwoman_episodes"
	case "Black Lightning":
		return "List_of_Black_Lightning_episodes"
	case "Constantine":
		return "List_of_Constantine_episodes"
	case "The Flash":
		return "List_of_The_Flash_(The_CW)_episodes"
	case "Freedom Fighters: The Ray":
		return "List_of_Freedom_Fighters:_The_Ray_episodes"
	case "DC's Legends of Tomorrow":
		return "List_of_DC%27s_Legends_of_Tomorrow_episodes"
	case "Stargirl":
		return "Stargirl_(TV_series)"
	case "Supergirl":
		return "List_of_Supergirl_episodes"
	case "Superman & Lois":
		return "List_of_Superman_%26_Lois_episodes"
	case "Vixen":
		return "List_of_Vixen_episodes"
	default:
		return slug
	}
}

func sourceSlug(name string) struct{ showSlug, root string } {
	return struct{ showSlug, root string }{
		showSlug: slugify(name),
		root:     "https://arrow.fandom.com/wiki/",
	}
}

// slugify turns an arbitrary series name into a stable kebab-case slug.
func slugify(s string) string {
	out := strings.ToLower(s)
	out = strings.ReplaceAll(out, " & ", "-and-")
	out = strings.ReplaceAll(out, " ", "-")
	out = strings.ReplaceAll(out, "'", "")
	return out
}

// urlCanonnicalise is a noop kept for parity with future URL-version
// harvesting. The signature returns the same URL string so the seed
// remains stable.
func urlCanonnicalise(in string) string {
	if u, err := url.Parse(in); err == nil {
		return u.String()
	}

	return in
}
