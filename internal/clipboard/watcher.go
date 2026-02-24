package clipboard

import "fmt"

// EventWatcher reports clipboard change events.
type EventWatcher interface {
	Events() <-chan struct{}
	Errors() <-chan error
	Close() error
}

// NewEventWatcher returns an event watcher for the current desktop session.
func NewEventWatcher() (EventWatcher, error) {
	switch sessionType() {
	case "wayland":
		if !hasCommand("wl-paste") {
			return nil, fmt.Errorf("wayland detected, but wl-paste not found")
		}
		return NewWaylandEventWatcher()
	case "x11":
		return NewX11EventWatcher()
	default:
		if hasCommand("wl-paste") {
			return NewWaylandEventWatcher()
		}
		return NewX11EventWatcher()
	}
}
