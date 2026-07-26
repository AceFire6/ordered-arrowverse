package handlers

import (
	"strings"
	"testing"
)

func TestStripJinja(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "removes a single-line jinja block, keeps body",
			in:   "<html>{% block x %}body{% endblock %}</html>",
			want: "<html>body</html>",
		},
		{
			name: "removes a multi-line jinja block, keeps body",
			in:   "<html>\n{% if foo %}\n<p>bar</p>\n{% endif %}\n</html>",
			want: "<html>\n\n<p>bar</p>\n\n</html>",
		},
		{
			name: "strips favicon expression as well as block",
			in:   `<head><link rel="icon" href="{{ static_url('favicon.png') }}" /></head>`,
			want: `<head><link rel="icon" href="" /></head>`,
		},
		{
			name: "leaves textual content alone",
			in:   "<p>Pure HTML with no Jinja.</p>",
			want: "<p>Pure HTML with no Jinja.</p>",
		},
		{
			name: "malformed block kept as-is",
			in:   "<html>{% unterminated</html>",
			want: "<html>{% unterminated</html>",
		},
		{
			name: "strips a generic expression too",
			in:   "<p>{{ value }}</p>",
			want: "<p></p>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripJinja(tt.in)
			if strings.TrimSpace(got) != strings.TrimSpace(tt.want) {
				t.Errorf("stripJinja mismatch\n got: %q\nwant: %q", got, tt.want)
			}
		})
	}
}

func TestBuildHomeFilterURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		hideList   string
		newestFirst bool
		want       string
	}{
		{
			name:     "no hide list",
			hideList: "",
			want:     "/",
		},
		{
			name:     "single hide",
			hideList: "arrow",
			want:     "/?hide_show=arrow",
		},
		{
			name:     "multiple hides",
			hideList: "arrow+flash+supergirl",
			want:     "/?hide_show=arrow&hide_show=flash&hide_show=supergirl",
		},
		{
			name:       "newest first with hide list",
			hideList:   "arrow+flash",
			newestFirst: true,
			want:       "/?hide_show=arrow&hide_show=flash&newest_first=true",
		},
		{
			name:       "newest first without hide list",
			hideList:   "",
			newestFirst: true,
			want:       "/?newest_first=true",
		},
		{
			name:     "empty hide values are dropped",
			hideList: "++flash",
			want:     "/?hide_show=flash",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildHomeFilterURL(tt.hideList, tt.newestFirst); got != tt.want {
				t.Errorf("buildHomeFilterURL(%q, %v) = %q, want %q",
					tt.hideList, tt.newestFirst, got, tt.want)
			}
		})
	}
}
