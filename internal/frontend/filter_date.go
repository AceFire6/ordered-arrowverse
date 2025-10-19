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
