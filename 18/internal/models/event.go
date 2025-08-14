package models

import "time"

// Event is an entry in the calendar
type Event struct {
	ID     int64     `json:"id"`
	UserID int64     `json:"user_id" validate:"required"`
	Date   time.Time `json:"date" validate:"required,datetime"`
	Text   string    `json:"event" validate:"required,min=1"`
}
