package clipboard

import (
	"bufio"
	"io"
	"os/exec"
)

// WaylandEventWatcher notifies when the Wayland clipboard changes.
type WaylandEventWatcher struct {
	cmd    *exec.Cmd
	events chan struct{}
	errs   chan error
	done   chan struct{}
}

// NewWaylandEventWatcher subscribes to Wayland clipboard change events.
func NewWaylandEventWatcher() (*WaylandEventWatcher, error) {
	cmd := exec.Command("wl-paste", "--watch", "sh", "-c", "printf '\\n'")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = io.Discard

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	w := &WaylandEventWatcher{
		cmd:    cmd,
		events: make(chan struct{}, 1),
		errs:   make(chan error, 1),
		done:   make(chan struct{}),
	}
	go w.loop(stdout)
	return w, nil
}

// Events returns a channel that receives on clipboard changes.
func (w *WaylandEventWatcher) Events() <-chan struct{} {
	return w.events
}

// Errors returns a channel that receives async errors.
func (w *WaylandEventWatcher) Errors() <-chan error {
	return w.errs
}

// Close stops the watcher process.
func (w *WaylandEventWatcher) Close() error {
	if w.cmd == nil || w.cmd.Process == nil {
		return nil
	}
	_ = w.cmd.Process.Kill()
	<-w.done
	w.cmd = nil
	return nil
}

func (w *WaylandEventWatcher) loop(r io.Reader) {
	defer close(w.done)
	defer close(w.events)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		select {
		case w.events <- struct{}{}:
		default:
		}
	}
	if err := scanner.Err(); err != nil {
		select {
		case w.errs <- err:
		default:
		}
	}
	if w.cmd != nil {
		if err := w.cmd.Wait(); err != nil {
			select {
			case w.errs <- err:
			default:
			}
		}
	}
}
