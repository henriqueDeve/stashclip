package clipboard

import (
	"bytes"
	"os/exec"
)

// WaylandBackend implements clipboard access using wl-copy/wl-paste.
type WaylandBackend struct{}

// NewWayland returns a Wayland clipboard backend.
func NewWayland() *WaylandBackend {
	return &WaylandBackend{}
}

// Read returns the current clipboard contents.
func (b *WaylandBackend) Read() (string, error) {
	cmd := exec.Command("wl-paste", "--no-newline")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Write updates the clipboard contents.
func (b *WaylandBackend) Write(text string) error {
	cmd := exec.Command("wl-copy")
	cmd.Stdin = bytes.NewBufferString(text)
	return cmd.Run()
}
