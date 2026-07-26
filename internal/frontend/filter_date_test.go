package frontend_test

import (
	"testing"
	"time"

	"github.com/AceFire6/ordered-arrowverse/internal/frontend"
)

func TestFilterDate_IsZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		seed func() *frontend.FilterDate
		want bool
	}{
		{
			name: "nil pointer",
			seed: func() *frontend.FilterDate { return nil },
			want: true,
		},
		{
			name: "zero value via address-of",
			seed: func() *frontend.FilterDate {
				fd := frontend.FilterDate(time.Time{})
				return &fd
			},
			want: true,
		},
		{
			name: "set value",
			seed: func() *frontend.FilterDate {
				d := time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)
				fd := frontend.FilterDate(d)
				return &fd
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.seed().IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterDate_TimePtr(t *testing.T) {
	t.Parallel()

	var nilFD *frontend.FilterDate
	if got := nilFD.TimePtr(); got != nil {
		t.Errorf("TimePtr() = %v, want nil", got)
	}

	zero := frontend.FilterDate(time.Time{})
	if got := (&zero).TimePtr(); got != nil {
		t.Errorf("TimePtr() on zero = %v, want nil", got)
	}

	want := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	fd := frontend.FilterDate(want)
	got := fd.TimePtr()
	if got == nil {
		t.Fatalf("TimePtr() returned nil for set date")
	}
	if !got.Equal(want) {
		t.Errorf("TimePtr() = %v, want %v", got, want)
	}
}

func TestFilterDate_String(t *testing.T) {
	t.Parallel()

	var nilFD *frontend.FilterDate
	if got := nilFD.String(); got != "" {
		t.Errorf("String() on nil = %q, want empty", got)
	}

	zero := frontend.FilterDate(time.Time{})
	if got := zero.String(); got != "" {
		t.Errorf("String() on zero value = %q, want empty", got)
	}

	want := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	fd := frontend.FilterDate(want)
	if got := fd.String(); got != "2024-01-02" {
		t.Errorf("String() = %q, want %q", got, "2024-01-02")
	}
}
