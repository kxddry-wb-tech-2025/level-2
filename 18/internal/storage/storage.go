package storage

import (
	"calendar/internal/models"
	"errors"
	"sync"
	"time"
)

var (
	// ErrEventNotFound is used when an event the user is trying to access is not found
	ErrEventNotFound = errors.New("event not found")

	// ErrUserNotFound is used when the user is not found
	ErrUserNotFound = errors.New("user not found")
)

// Storage is a struct suitable for a calendar.
type Storage struct {
	byUser map[int64][]*models.Event
	mp     map[int64]*models.Event
	idx    map[int64]int
	nextID int64
	mu     *sync.RWMutex
}

// New initializes a new Storage
func New() *Storage {
	return &Storage{
		byUser: make(map[int64][]*models.Event),
		mp:     make(map[int64]*models.Event),
		idx:    make(map[int64]int),
		nextID: 1,
		mu:     new(sync.RWMutex),
	}
}

// Create creates a new event
func (s *Storage) Create(userID int64, date time.Time, text string) *models.Event {
	s.mu.Lock()
	defer s.mu.Unlock()

	e := &models.Event{
		UserID: userID,
		Date:   date,
		Text:   text,
		ID:     s.nextID,
	}

	s.mp[s.nextID] = e
	s.idx[s.nextID] = len(s.byUser)
	s.byUser[userID] = append(s.byUser[userID], e)
	s.nextID++
	return e
}

// Update updates an existing event and returns ErrEventNotFound if it wasn't found
func (s *Storage) Update(id int64, date time.Time, text string) (*models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.mp[id]
	if !ok {
		return nil, ErrEventNotFound
	}

	e.Date = date
	e.Text = text
	uid := e.UserID
	idx := s.idx[id]
	s.byUser[uid][idx] = e
	s.mp[id] = e

	return e, nil
}

// Delete deletes an event and returns ErrEventNotFound if it was not found
func (s *Storage) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.mp[id]
	if !ok {
		return ErrEventNotFound
	}
	uid := e.UserID
	idx := s.idx[id]
	last := len(s.byUser[uid]) - 1
	lastID := s.byUser[uid][last].ID

	delete(s.mp, id)
	delete(s.idx, id)

	s.byUser[uid][idx], s.byUser[uid][last] = s.byUser[uid][last], s.byUser[uid][idx]
	s.byUser[uid] = s.byUser[uid][:last]
	s.idx[lastID] = idx

	return nil
}

func (s *Storage) withFilter(userID int64, filter func(e *models.Event) bool) ([]*models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.byUser[userID]) == 0 {
		return nil, ErrUserNotFound
	}

	var out []*models.Event
	for _, e := range s.byUser[userID] {
		if filter(e) {
			out = append(out, e)
		}
	}
	return out, nil
}

// GetDay shows the events with the same day as requested
func (s *Storage) GetDay(userID int64) ([]*models.Event, error) {
	y, m, d := time.Now().Date()
	filter := func(e *models.Event) bool {
		yy, mm, dd := e.Date.Date()
		diff := time.Date(yy, mm, dd, 0, 0, 0, 0, time.UTC).Sub(time.Date(y, m, d, 0, 0, 0, 0, time.UTC))
		return diff <= time.Hour*24
	}
	return s.withFilter(userID, filter)
}

// GetMonth shows the events following the requested date within a month
func (s *Storage) GetMonth(userID int64) ([]*models.Event, error) {
	y, m, d := time.Now().Date()
	filter := func(e *models.Event) bool {
		yy, mm, dd := e.Date.Date()
		diff := time.Date(yy, mm, dd, 0, 0, 0, 0, time.UTC).Sub(time.Date(y, m, d, 0, 0, 0, 0, time.UTC))
		return diff <= time.Hour*24*30
	}
	return s.withFilter(userID, filter)
}

// GetWeek shows the events following the requested date within a week
func (s *Storage) GetWeek(userID int64) ([]*models.Event, error) {
	y, m, d := time.Now().Date()
	filter := func(e *models.Event) bool {
		yy, mm, dd := e.Date.Date()
		diff := time.Date(yy, mm, dd, 0, 0, 0, 0, time.UTC).Sub(time.Date(y, m, d, 0, 0, 0, 0, time.UTC))
		return diff <= time.Hour*24*7
	}
	return s.withFilter(userID, filter)
}
