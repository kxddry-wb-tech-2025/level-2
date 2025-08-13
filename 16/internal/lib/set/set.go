package set

import "sync"

// Set is a concurrency-safe data structure like the set in Python
type Set struct {
	data map[any]struct{}
	mu   *sync.RWMutex
}

// New creates a set
func New() *Set {
	return &Set{
		data: make(map[any]struct{}),
		mu:   new(sync.RWMutex),
	}
}

// Add adds an item to a set
func (s *Set) Add(item any) {
	s.mu.Lock()
	s.data[item] = struct{}{}
	s.mu.Unlock()
}

// Remove removes an item from the set
func (s *Set) Remove(item any) {
	s.mu.Lock()
	delete(s.data, item)
	s.mu.Unlock()
}

// Contains tells us whether the set contains an element
func (s *Set) Contains(item any) bool {
	s.mu.RLock()
	_, ok := s.data[item]
	s.mu.RUnlock()
	return ok
}

// Len tells us the length of the set
func (s *Set) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}
