package frontend

type ShowData struct {
	Name    string
	URL     string
	RootURL string
}

type TableRow struct {
	ShowSlug    string
	Series      string
	EpisodeId   string
	EpisodeName string
	AirDate     string
	SourceLink  string
	RowNumber   int
}
