package storage

import "sync"

// Storage represents a thread-safe in-memory storage.
type Storage struct {
	mu    sync.RWMutex
	store map[string]string
}

// NewStorage creates a new Storage instance.
func NewStorage() *Storage {
	return &Storage{
		store: make(map[string]string),
	}
}

// Set stores a key-value pair in the storage.
func (s *Storage) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[key] = value
}

// Get retrieves a value by key from the storage.
func (s *Storage) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, exists := s.store[key]
	return value, exists
}

// Store returns a copy of the current storage.
func (s *Storage) Store() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	persistenceCopy := make(map[string]string)
	for key, value := range s.store {
		persistenceCopy[key] = value
	}
	return persistenceCopy
}
