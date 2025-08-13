package storage

import (
	"calendar/internal/models"
	"errors"
	"sync"
	"time"
)

var (
	ErrEventExists   = errors.New("event already exists")
	ErrEventNotFound = errors.New("event not found")
	ErrUserNotFound  = errors.New("user not found")
)

type Storage struct {
	byUser map[int64][]*models.Event
	mp     map[int64]*models.Event
	idx    map[int64]int
	nextID int64
	mu     *sync.RWMutex
}

func New() *Storage {
	return &Storage{
		byUser: make(map[int64][]*models.Event),
		mp:     make(map[int64]*models.Event),
		idx:    make(map[int64]int),
		nextID: 1,
		mu:     new(sync.RWMutex),
	}
}

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

func (s *Storage) Update(id int64, date time.Time, text string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.mp[id]
	if !ok {
		return ErrEventNotFound
	}

	e.Date = date
	e.Text = text
	uid := e.UserID
	idx := s.idx[id]
	s.byUser[uid][idx] = e
	s.mp[id] = e

	return nil
}

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
	s.byUser[uid][last], s.byUser[uid][idx] = s.byUser[uid][idx], nil
	s.idx[s.byUser[uid][last].ID] = idx
	s.byUser[uid] = s.byUser[uid][:last]

	delete(s.idx, id)
	delete(s.mp, id)
	return nil
}

func (s *Storage) GetDay(userID int64, date time.Time) ([]*models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.byUser[userID]) == 0 {
		return nil, ErrUserNotFound
	}
	var out []*models.Event
	for _, e := range s.byUser[userID] {
		if e.Date.Equal(date) {
			out = append(out, e)
		}
	}
	return out, nil

}
