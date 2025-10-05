package frontend_test

import (
	"reflect"
	"testing"

	"github.com/AceFire6/ordered-arrowverse/internal/frontend"
)

func TestMergePageConfigs(t *testing.T) {
	pageConfigA := frontend.PageConfig{
		Title:       "A",
		Heading:     "A Heading",
		Reverser:    fakeReverser{Name: "A"},
		Contents:    nil,
		JsScripts:   []string{"script-a1.js", "script-a2.js"},
		Stylesheets: []string{"style-a1.css"},
	}

	pageConfigB := frontend.PageConfig{
		Title:       "B",
		Heading:     "B Heading",
		Reverser:    fakeReverser{Name: "B"},
		Contents:    nil,
		JsScripts:   []string{"script-b1.js"},
		Stylesheets: []string{"style-b1.css", "style-b2.css"},
	}

	type args struct {
		pageConfigA *frontend.PageConfig
		pageConfigB *frontend.PageConfig
		option      frontend.MergeOption
	}
	tests := []struct {
		want *frontend.PageConfig
		args args
		name string
	}{
		{
			name: "MergeAll - merges config A and B correctly - overwrites A with B's values and merges JS & CSS",
			args: args{
				pageConfigA: pageConfigA.Copy(),
				pageConfigB: pageConfigB.Copy(),
				option:      frontend.MergeAll,
			},
			want: &frontend.PageConfig{
				Title:       pageConfigB.Title,
				Heading:     pageConfigB.Heading,
				Reverser:    pageConfigB.Reverser,
				Contents:    pageConfigB.Contents,
				JsScripts:   []string{"script-a1.js", "script-a2.js", "script-b1.js"},
				Stylesheets: []string{"style-a1.css", "style-b1.css", "style-b2.css"},
			},
		},
		{
			name: "MergeAll - merges config A and B correctly - does not overwrite A.Title if B.Title is empty",
			args: args{
				pageConfigA: pageConfigA.Copy(),
				pageConfigB: &frontend.PageConfig{
					Title:       "",
					Heading:     "B Heading",
					Reverser:    fakeReverser{Name: "B"},
					Contents:    nil,
					JsScripts:   []string{"script-b1.js"},
					Stylesheets: []string{"style-b1.css", "style-b2.css"},
				},
				option: frontend.MergeAll,
			},
			want: &frontend.PageConfig{
				Title:       pageConfigA.Title,
				Heading:     "B Heading",
				Reverser:    pageConfigB.Reverser,
				Contents:    pageConfigB.Contents,
				JsScripts:   []string{"script-a1.js", "script-a2.js", "script-b1.js"},
				Stylesheets: []string{"style-a1.css", "style-b1.css", "style-b2.css"},
			},
		},
		{
			name: "MergeAll - merges config A and B correctly - does not overwrite A.Heading if B.Heading is empty",
			args: args{
				pageConfigA: pageConfigA.Copy(),
				pageConfigB: &frontend.PageConfig{
					Title:       "B Title",
					Heading:     "",
					Reverser:    fakeReverser{Name: "B"},
					Contents:    nil,
					JsScripts:   []string{"script-b1.js"},
					Stylesheets: []string{"style-b1.css", "style-b2.css"},
				},
				option: frontend.MergeAll,
			},
			want: &frontend.PageConfig{
				Title:       "B Title",
				Heading:     pageConfigA.Heading,
				Reverser:    pageConfigB.Reverser,
				Contents:    pageConfigB.Contents,
				JsScripts:   []string{"script-a1.js", "script-a2.js", "script-b1.js"},
				Stylesheets: []string{"style-a1.css", "style-b1.css", "style-b2.css"},
			},
		},
		{
			name: "OverwriteAll - merges config A and B correctly - overwrites A with B's values",
			args: args{
				pageConfigA: pageConfigA.Copy(),
				pageConfigB: pageConfigB.Copy(),
				option:      frontend.OverwriteAll,
			},
			want: pageConfigB.Copy(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := frontend.MergePageConfigs(tt.args.pageConfigA, tt.args.pageConfigB, tt.args.option); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergePageConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}

type fakeReverser struct {
	Name string
}

func (fr fakeReverser) Reverse(name string, params ...interface{}) string {
	return fr.Name
}
