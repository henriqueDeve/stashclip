package clipboard

import (
	"os"
	"os/exec"
	"strings"
)

func sessionType() string {
	t := strings.ToLower(strings.TrimSpace(os.Getenv("XDG_SESSION_TYPE")))
	if t == "wayland" || t == "x11" {
		return t
	}
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return "wayland"
	}
	if os.Getenv("DISPLAY") != "" {
		return "x11"
	}
	return ""
}

func hasCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
