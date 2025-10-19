package db

import (
	"fmt"

	"github.com/AceFire6/ordered-arrowverse/internal/frontend"
)

func ShowListRowToShowData(slr GetShowListRow) frontend.ShowData {
	copiedDataSources := make([]string, len(slr.DataSources))
	copy(copiedDataSources, slr.DataSources)

	return frontend.ShowData{
		Slug:        slr.Slug,
		Name:        slr.Name,
		DataSources: copiedDataSources,
	}
}

func ShowListRowsToShowData(showListRows []GetShowListRow) []frontend.ShowData {
	showList := make([]frontend.ShowData, len(showListRows))
	for idx, dbShow := range showListRows {
		showList[idx] = ShowListRowToShowData(dbShow)
	}

	return showList
}

func EpisodeRowsToTableRow(episodeRows []GetEpisodesRow) []frontend.TableRow {
	tableRows := make([]frontend.TableRow, len(episodeRows))
	for i, episode := range episodeRows {
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

	return tableRows
}
