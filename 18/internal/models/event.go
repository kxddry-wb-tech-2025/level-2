package models

import "time"

type Event struct {
	ID     int64
	UserID int64
	Date   time.Time `json:"date"`
	Text   string    `json:"event"`
}
