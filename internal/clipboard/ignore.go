package clipboard

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"stashclip/internal/store"
)

type ignoredEntry struct {
	Hash      string    `json:"hash"`
	ExpiresAt time.Time `json:"expires_at"`
}

const ignoredTTL = 10 * time.Second

// MarkIgnored marks clipboard text as app-originated so daemon can skip storing it once.
func MarkIgnored(text string) error {
	entry := ignoredEntry{
		Hash:      hashText(text),
		ExpiresAt: time.Now().Add(ignoredTTL),
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	path := ignoredPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// ShouldIgnore reports whether clipboard text should be ignored by storage capture.
func ShouldIgnore(text string) bool {
	path := ignoredPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	var entry ignoredEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		_ = os.Remove(path)
		return false
	}
	if time.Now().After(entry.ExpiresAt) {
		_ = os.Remove(path)
		return false
	}
	if entry.Hash != hashText(text) {
		return false
	}
	_ = os.Remove(path)
	return true
}

func ignoredPath() string {
	base := store.DefaultPath()
	if base == "" {
		return filepath.Join("/tmp", "stashclip-ignore.json")
	}
	dir := filepath.Dir(base)
	if err := os.MkdirAll(dir, 0o755); err != nil && !errors.Is(err, os.ErrExist) {
		return filepath.Join("/tmp", "stashclip-ignore.json")
	}
	return filepath.Join(dir, "ignore.json")
}

func hashText(text string) string {
	sum := sha256.Sum256([]byte(text))
	return hex.EncodeToString(sum[:])
}
