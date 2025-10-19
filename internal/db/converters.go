package db

import (
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
