package store

import (
	"sync"
	"time"
)

// Entry represents a stored clipboard item.
type Entry struct {
	Text     string
	AddedAt  time.Time
}

// Store keeps clipboard entries in memory.
type Store struct {
	mu      sync.Mutex
	entries []Entry
}

// New returns an empty store.
func New() *Store {
	return &Store{}
}

// Add inserts a new entry unless it is a consecutive duplicate.
func (s *Store) Add(text string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.entries) > 0 && s.entries[len(s.entries)-1].Text == text {
		return
	}

	s.entries = append(s.entries, Entry{Text: text, AddedAt: time.Now()})
	if len(s.entries) > 200 {
		s.entries = s.entries[len(s.entries)-200:]
	}
}

// List returns a copy of all entries.
func (s *Store) List() []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]Entry, len(s.entries))
	copy(out, s.entries)
	return out
}

// Clear removes all entries.
func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries = nil
}
