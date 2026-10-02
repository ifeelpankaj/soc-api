package models

import (
	"encoding/json"
	"time"
)

// DateOnly serializes a calendar date as YYYY-MM-DD in JSON responses.
type DateOnly time.Time

func NewDateOnly(value *time.Time) *DateOnly {
	if value == nil {
		return nil
	}

	d := DateOnly(*value)
	return &d
}

func (d DateOnly) Time() time.Time {
	return time.Time(d)
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	if t.IsZero() {
		return []byte("null"), nil
	}

	return json.Marshal(t.UTC().Format("2006-01-02"))
}
