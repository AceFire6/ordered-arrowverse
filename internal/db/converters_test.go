package db_test

import (
	"testing"
	"time"

	"github.com/AceFire6/ordered-arrowverse/internal/db"
	"github.com/AceFire6/ordered-arrowverse/internal/frontend"
)

func TestEpisodeRowsToTableRow(t *testing.T) {
	t.Parallel()

	airDate := time.Date(2020, 5, 12, 0, 0, 0, 0, time.UTC)
	rows := []db.GetEpisodesRow{
		{
			ShowName:       "Arrow",
			ShowSlug:       "arrow",
			EpisodeID:      1,
			Name:           "Pilot",
			AirDate:        airDate,
			Season:         1,
			Episode:        1,
			Description:    "Pilot episode",
			SourceLink:     "https://arrow.fandom.com/wiki/Pilot_(Arrow)",
			Ordering:       1,
			ShowID:         1,
			CreatedAt:      airDate,
			ModifiedAt:     airDate,
			RuntimeMinutes: 42,
		},
	}

	got := db.EpisodeRowsToTableRow(rows)
	if len(got) != 1 {
		t.Fatalf("expected 1 row, got %d", len(got))
	}
	row := got[0]
	if row.ShowSlug != "arrow" {
		t.Errorf("ShowSlug = %q, want arrow", row.ShowSlug)
	}
	if row.RowNumber != 1 {
		t.Errorf("RowNumber = %d, want 1 (post-filter 1-based)", row.RowNumber)
	}
	if row.Series != "Arrow" {
		t.Errorf("Series = %q, want Arrow", row.Series)
	}
	if row.EpisodeId != "S01E01" {
		t.Errorf("EpisodeId = %q, want S01E01", row.EpisodeId)
	}
	if row.EpisodeName != "Pilot" {
		t.Errorf("EpisodeName = %q, want Pilot", row.EpisodeName)
	}
	if row.AirDate != "May 12, 2020" {
		t.Errorf("AirDate = %q, want May 12, 2020", row.AirDate)
	}
}

func TestShowListRowsToShowData_DefensiveCopy(t *testing.T) {
	t.Parallel()

	row := db.GetShowListRow{
		Slug:        "arrow",
		Name:        "Arrow",
		DataSources: []string{"https://arrow.fandom.com/wiki/List_of_Arrow_episodes"},
	}

	got := db.ShowListRowsToShowData([]db.GetShowListRow{row})
	if len(got) != 1 {
		t.Fatalf("expected 1 show, got %d", len(got))
	}
	if got[0].Slug != "arrow" {
		t.Errorf("Slug = %q, want arrow", got[0].Slug)
	}
	if len(got[0].DataSources) != 1 || got[0].DataSources[0] != row.DataSources[0] {
		t.Errorf("DataSources mismatch")
	}

	// Confirm ShowListRowToShowData performs a defensive copy so the
	// caller can mutate the result without touching the input.
	row.DataSources[0] = "mutated"
	if got[0].DataSources[0] == "mutated" {
		t.Errorf("ShowListRowsToShowData did not copy DataSources")
	}

	// Modifying the returned slice should not affect the input row.
	got[0].DataSources[0] = "modified-output"
	if row.DataSources[0] == "modified-output" {
		t.Errorf("returned DataSources shares storage with input")
	}

	// Sanity: the result implements the ShowData interface used by templ.
	var _ frontend.ShowData = got[0]
}
