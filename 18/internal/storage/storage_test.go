package storage_test

import (
	"calendar/internal/models"
	"calendar/internal/storage"
	"errors"
	"testing"
	"time"
)

func TestCreate_HappyPath(t *testing.T) {
	s := storage.New()
	date := models.Date(time.Now())

	e1 := s.Create(1, date, "Event 1")
	e2 := s.Create(1, date, "Event 2") // same date, same user

	if e1.ID == e2.ID {
		t.Errorf("expected different IDs for two events, got %d and %d", e1.ID, e2.ID)
	}

	events, _ := s.GetDay(1, date)
	found := 0
	for _, e := range events {
		if e.ID == e1.ID || e.ID == e2.ID {
			found++
		}
	}
	if found != 2 {
		t.Errorf("expected 2 events for today, got %d", found)
	}
}

func TestCreate_NonHappyPath(t *testing.T) {
	// Creating same event is allowed, no non-happy path for Create
}

func TestUpdate_HappyPath(t *testing.T) {
	s := storage.New()
	date := time.Now()
	e := s.Create(1, models.Date(date), "Initial")

	newDate := date.Add(24 * time.Hour)
	_, err := s.Update(e.ID, models.Date(newDate), "Updated")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events, _ := s.GetDay(1, models.Date(newDate))
	found := false
	for _, ev := range events {
		if ev.ID == e.ID && ev.Text == "Updated" {
			found = true
		}
	}
	if !found {
		t.Errorf("update failed, event not found for today/week/month: %v", events)
	}

	// Updating with same values
	_, err = s.Update(e.ID, models.Date(newDate), "Updated")
	if err != nil {
		t.Errorf("updating with same values should succeed, got %v", err)
	}
}

func TestUpdate_NonHappyPath(t *testing.T) {
	s := storage.New()
	_, err := s.Update(999, models.Date(time.Now()), "Fail")
	if !errors.Is(err, storage.ErrEventNotFound) {
		t.Errorf("expected ErrEventNotFound, got %v", err)
	}
}

func TestDelete_HappyPath(t *testing.T) {
	s := storage.New()
	date := time.Now()
	e1 := s.Create(1, models.Date(date), "Event 1")
	e2 := s.Create(1, models.Date(date), "Event 2")

	err := s.Delete(e1.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events, _ := s.GetDay(1, models.Date(date))
	for _, ev := range events {
		if ev.ID == e1.ID {
			t.Errorf("delete failed, event still present: %v", ev)
		}
	}

	// Delete last element
	err = s.Delete(e2.ID)
	if err != nil {
		t.Fatalf("unexpected error when deleting last element: %v", err)
	}
	events, _ = s.GetDay(1, models.Date(date))
	if len(events) != 0 {
		t.Errorf("expected 0 events after deleting all, got %d", len(events))
	}
}

func TestDelete_NonHappyPath(t *testing.T) {
	s := storage.New()
	err := s.Delete(999)
	if !errors.Is(err, storage.ErrEventNotFound) {
		t.Errorf("expected ErrEventNotFound, got %v", err)
	}
}

func TestGetDay_HappyPath(t *testing.T) {
	s := storage.New()
	date := time.Now()
	s.Create(1, models.Date(date), "Day event")

	events, err := s.GetDay(1, models.Date(date))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) == 0 {
		t.Errorf("expected at least 1 event for today, got 0")
	}
}

func TestGetDay_NonHappyPath(t *testing.T) {
	s := storage.New()
	date := time.Now()
	_, err := s.GetDay(999, models.Date(date))
	if !errors.Is(err, storage.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestGetWeek_HappyPath(t *testing.T) {
	s := storage.New()
	base := time.Now()
	s.Create(1, models.Date(base), "Week event")
	s.Create(1, models.Date(base.Add(6*24*time.Hour)), "Week end event")

	events, err := s.GetWeek(1, models.Date(base))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := 0
	for _, e := range events {
		if e.Text == "Week event" || e.Text == "Week end event" {
			found++
		}
	}
	if found != 2 {
		t.Errorf("expected 2 events for this week, got %d", found)
	}
}

func TestGetWeek_NonHappyPath(t *testing.T) {
	s := storage.New()
	_, err := s.GetWeek(999, models.Date(time.Now()))
	if !errors.Is(err, storage.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestGetMonth_HappyPath(t *testing.T) {
	s := storage.New()
	base := time.Now()
	s.Create(1, models.Date(base), "Month event")
	s.Create(1, models.Date(base.Add(29*24*time.Hour)), "End of month event")

	events, err := s.GetMonth(1, models.Date(base))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := 0
	for _, e := range events {
		if e.Text == "Month event" || e.Text == "End of month event" {
			found++
		}
	}
	if found != 2 {
		t.Errorf("expected 2 events for this month, got %d", found)
	}
}

func TestGetMonth_NonHappyPath(t *testing.T) {
	s := storage.New()
	_, err := s.GetMonth(999, models.Date(time.Now()))
	if !errors.Is(err, storage.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}
