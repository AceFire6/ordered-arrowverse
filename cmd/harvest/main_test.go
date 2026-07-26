package main

import (
	"testing"
)

func TestTruncateDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{"2020-05-12T00:00:00", "2020-05-12"},
		{"2024-01-02T15:04:05Z", "2024-01-02"},
		{"2024-01-02", "2024-01-02"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := truncateDate(tt.in); got != tt.want {
				t.Errorf("truncateDate(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestSQLString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{"Arrow", "'Arrow'"},
		{"DC's Legends", "'DC''s Legends'"},
		{"", "''"},
		{"already'has'quotes", "'already''has''quotes'"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := sqlString(tt.in); got != tt.want {
				t.Errorf("sqlString(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{"Arrow", "arrow"},
		{"Black Lightning", "black-lightning"},
		{"DC's Legends of Tomorrow", "dcs-legends-of-tomorrow"},
		{"Superman & Lois", "superman-and-lois"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := slugify(tt.in); got != tt.want {
				t.Errorf("slugify(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestShowSlugRelativePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		series  string
		wantSub string
	}{
		{"Arrow", "Arrow", "List_of_Arrow_episodes"},
		{"Legends", "DC's Legends of Tomorrow", "List_of_DC%27s_Legends_of_Tomorrow_episodes"},
		{"Stargirl", "Stargirl", "Stargirl_(TV_series)"},
		{"Superman & Lois", "Superman & Lois", "List_of_Superman_%26_Lois_episodes"},
		{"Unknown", "Some Random Show", "some-random-show"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := showSlugRelativePath(tt.series, slugify(tt.series))
			if !contains(got, tt.wantSub) {
				t.Errorf("showSlugRelativePath(%q) = %q, expected to contain %q",
					tt.series, got, tt.wantSub)
			}
		})
	}
}

func contains(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
