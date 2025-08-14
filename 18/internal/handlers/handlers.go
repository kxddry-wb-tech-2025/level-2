package handlers

import (
	"calendar/internal/models"
	"time"
)

// Storage is an interface required for the calendar.
type Storage interface {
	Create(userID int64, date time.Time, text string) *models.Event
	Update(id int64, date time.Time, text string) error
	Delete(id int64) error
	GetDay(userID int64, date time.Time) ([]*models.Event, error)
	GetWeek(userID int64, date time.Time) ([]*models.Event, error)
	GetMonth(userID int64, date time.Time) ([]*models.Event, error)
}
