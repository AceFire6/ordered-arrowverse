package frontend

type ShowData struct {
	Slug        string
	Name        string
	DataSources []string
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
