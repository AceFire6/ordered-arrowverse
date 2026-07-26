package frontend

import (
	"time"

	"github.com/goccy/go-json"
)

type FilterDate time.Time

func (fd *FilterDate) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}

	t, err := time.Parse("2006-01-02", string(text))
	if err != nil {
		return err
	}

	*fd = FilterDate(t)
	return nil
}

func (fd *FilterDate) String() string {
	// Return an empty string when the filter date is an empty value
	if fd == nil {
		return ""
	}

	return time.Time(*fd).Format("2006-01-02")
}

func (fd FilterDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(fd.String())
}

// IsZero reports whether the FilterDate is nil or holds a zero time.Time.
// A zero FilterDate is treated as "no filter applied".
func (fd *FilterDate) IsZero() bool {
	if fd == nil {
		return true
	}

	return time.Time(*fd).IsZero()
}

// TimePtr returns a pointer to the underlying time.Time, or nil when the
// filter is unset. Useful for handing off to a DB layer that distinguishes
// nil from zero via *time.Time.
func (fd *FilterDate) TimePtr() *time.Time {
	if fd.IsZero() {
		return nil
	}

	t := time.Time(*fd)

	return &t
}
