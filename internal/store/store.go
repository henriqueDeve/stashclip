package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Entry represents a stored clipboard item.
type Entry struct {
	Text    string
	AddedAt time.Time
}

// Store keeps clipboard entries in memory.
type Store struct {
	mu      sync.Mutex
	entries []Entry
	path    string
}

// New returns a store backed by the default on-disk path.
func New() (*Store, error) {
	return NewWithPath(DefaultPath())
}

// NewWithPath returns a store backed by a specific on-disk path.
func NewWithPath(path string) (*Store, error) {
	s := &Store{path: path}
	if path == "" {
		return s, nil
	}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

// DefaultPath returns the default storage location.
func DefaultPath() string {
	if base := os.Getenv("XDG_DATA_HOME"); base != "" {
		return filepath.Join(base, "stashclip", "store.json")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".local", "share", "stashclip", "store.json")
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
	_ = s.save()
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
	_ = s.save()
}

func (s *Store) load() error {
	if s.path == "" {
		return nil
	}
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}
	s.entries = entries
	return nil
}

func (s *Store) save() error {
	if s.path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(s.entries)
	if err != nil {
		return err
	}
	tmpPath := s.path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpPath, s.path)
}
