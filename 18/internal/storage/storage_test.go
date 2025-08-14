package storage_test

import (
	"calendar/internal/storage"
	"errors"
	"testing"
	"time"
)

func TestCreate_HappyPath(t *testing.T) {
	s := storage.New()
	date := time.Now()

	e1 := s.Create(1, date, "Event 1")
	e2 := s.Create(1, date, "Event 2") // same date, same user

	if e1.ID == e2.ID {
		t.Errorf("expected different IDs for two events, got %d and %d", e1.ID, e2.ID)
	}

	events, _ := s.GetDay(1, date)
	if len(events) != 2 {
		t.Errorf("expected 2 events for the same day, got %d", len(events))
	}
}

func TestCreate_NonHappyPath(t *testing.T) {
	// In this storage, creating the same event is allowed
	// so no non-happy path for Create itself
}

func TestUpdate_HappyPath(t *testing.T) {
	s := storage.New()
	date := time.Now()
	e := s.Create(1, date, "Initial")

	newDate := date.Add(time.Hour * 24)
	_, err := s.Update(e.ID, newDate, "Updated")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events, _ := s.GetDay(1, newDate)
	if len(events) != 1 || events[0].Text != "Updated" {
		t.Errorf("update failed, got %v", events)
	}

	// Updating with same values
	_, err = s.Update(e.ID, newDate, "Updated")
	if err != nil {
		t.Errorf("updating with same values should succeed, got %v", err)
	}
}

func TestUpdate_NonHappyPath(t *testing.T) {
	s := storage.New()
	_, err := s.Update(999, time.Now(), "Fail")
	if !errors.Is(err, storage.ErrEventNotFound) {
		t.Errorf("expected ErrEventNotFound, got %v", err)
	}
}

func TestDelete_HappyPath(t *testing.T) {
	s := storage.New()
	date := time.Now()
	e1 := s.Create(1, date, "Event 1")
	e2 := s.Create(1, date, "Event 2")

	err := s.Delete(e1.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events, _ := s.GetDay(1, date)
	if len(events) != 1 || events[0].ID != e2.ID {
		t.Errorf("delete failed, remaining events: %v", events)
	}

	// Delete last element
	err = s.Delete(e2.ID)
	if err != nil {
		t.Fatalf("unexpected error when deleting last element: %v", err)
	}
	events, _ = s.GetDay(1, date)
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
	date := time.Date(2025, 8, 14, 10, 0, 0, 0, time.UTC)
	s.Create(1, date, "Day event")

	events, err := s.GetDay(1, date)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}

func TestGetDay_NonHappyPath(t *testing.T) {
	s := storage.New()
	_, err := s.GetDay(999, time.Now())
	if !errors.Is(err, storage.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestGetWeek_HappyPath(t *testing.T) {
	s := storage.New()
	base := time.Date(2025, 8, 14, 10, 0, 0, 0, time.UTC)
	s.Create(1, base, "Week event")
	s.Create(1, base.Add(6*24*time.Hour), "Week end event")

	events, err := s.GetWeek(1, base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
}

func TestGetWeek_NonHappyPath(t *testing.T) {
	s := storage.New()
	_, err := s.GetWeek(999, time.Now())
	if !errors.Is(err, storage.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestGetMonth_HappyPath(t *testing.T) {
	s := storage.New()
	base := time.Date(2025, 8, 1, 10, 0, 0, 0, time.UTC)
	s.Create(1, base, "Month event")
	s.Create(1, base.Add(29*24*time.Hour), "End of month event")

	events, err := s.GetMonth(1, base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
}

func TestGetMonth_NonHappyPath(t *testing.T) {
	s := storage.New()
	_, err := s.GetMonth(999, time.Now())
	if !errors.Is(err, storage.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}
