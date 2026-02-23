package clipboard

import (
	"bytes"
	"os/exec"
)

// X11Backend implements clipboard access using xclip.
type X11Backend struct{}

// NewX11 returns an X11 clipboard backend.
func NewX11() *X11Backend {
	return &X11Backend{}
}

// Read returns the current clipboard contents.
func (b *X11Backend) Read() (string, error) {
	return clipRead()
}

// Write updates the clipboard contents.
func (b *X11Backend) Write(text string) error {
	return clipWrite(text)
}

func clipRead() (string, error) {
	cmd := exec.Command("xclip", "-selection", "clipboard", "-o")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func clipWrite(text string) error {
	cmd := exec.Command("xclip", "-selection", "clipboard", "-i")
	cmd.Stdin = bytes.NewBufferString(text)
	return cmd.Run()
}
