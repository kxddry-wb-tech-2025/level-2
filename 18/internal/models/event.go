package models

import (
	"strings"
	"time"
)

// Event is an entry in the calendar
type Event struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id" validate:"required"`
	Date   Date   `json:"date" validate:"required,datetime"`
	Text   string `json:"event" validate:"required,min=1"`
}

type Date time.Time

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	*d = Date(t)
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(d).Format("2006-01-02") + `"`), nil
}
