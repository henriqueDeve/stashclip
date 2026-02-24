package daemon

import (
	"crypto/sha256"
	"os"
	"os/signal"
	"syscall"

	"stashclip/internal/clipboard"
	"stashclip/internal/store"
)

// Run starts the clipboard monitoring loop and blocks until interrupted.
func Run(clipboardProvider clipboard.ClipboardProvider, store *store.Store) error {
	watcher, err := clipboard.NewEventWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sig)

	var lastHash [32]byte
	var hasHash bool

	for {
		select {
		case <-sig:
			return nil
		case err := <-watcher.Errors():
			if err != nil {
				return err
			}
		case _, ok := <-watcher.Events():
			if !ok {
				return nil
			}
			text, err := clipboardProvider.Read()
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
