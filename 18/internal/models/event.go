package models

import "time"

// Event is an entry in the calendar
type Event struct {
	ID     int64
	UserID int64
	Date   time.Time `json:"date"`
	Text   string    `json:"event"`
}
