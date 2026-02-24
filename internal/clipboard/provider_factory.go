package clipboard

import "fmt"

// NewProvider returns a clipboard provider for the current desktop session.
func NewProvider() (ClipboardProvider, error) {
	switch sessionType() {
	case "wayland":
		if hasCommand("wl-copy") && hasCommand("wl-paste") {
			return NewWayland(), nil
		}
		return nil, fmt.Errorf("wayland detected, but wl-copy/wl-paste not found")
	case "x11":
		if hasCommand("xclip") {
			return NewX11(), nil
		}
		return nil, fmt.Errorf("x11 detected, but xclip not found")
	default:
		if hasCommand("wl-copy") && hasCommand("wl-paste") {
			return NewWayland(), nil
		}
		if hasCommand("xclip") {
			return NewX11(), nil
		}
		return nil, fmt.Errorf("no clipboard backend found (need wl-copy/wl-paste or xclip)")
	}
}
