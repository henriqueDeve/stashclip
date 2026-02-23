package daemon

import (
	"crypto/sha256"
	"os"
	"os/signal"
	"time"

	"stashclip/internal/clipboard"
	"stashclip/internal/store"
)

// Run starts the clipboard monitoring loop and blocks until interrupted.
func Run(backend clipboard.Backend, store *store.Store) error {
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	defer signal.Stop(sig)

	var lastHash [32]byte
	var hasHash bool

	for {
		select {
		case <-sig:
			return nil
		case <-ticker.C:
			text, err := backend.Read()
			if err != nil {
				continue
			}
			hash := sha256.Sum256([]byte(text))
			if hasHash && hash == lastHash {
				continue
			}
			hasHash = true
			lastHash = hash
			store.Add(text)
		}
	}
}
