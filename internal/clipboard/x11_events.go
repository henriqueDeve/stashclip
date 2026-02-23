package clipboard

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xfixes"
	"github.com/BurntSushi/xgb/xproto"
)

// X11EventWatcher notifies when the X11 clipboard changes.
type X11EventWatcher struct {
	conn   *xgb.Conn
	events chan struct{}
	errs   chan error
}

// NewX11EventWatcher subscribes to X11 clipboard change events.
func NewX11EventWatcher() (*X11EventWatcher, error) {
	conn, err := xgb.NewConn()
	if err != nil {
		return nil, err
	}

	if err := xfixes.Init(conn); err != nil {
		conn.Close()
		return nil, err
	}

	root := xproto.Setup(conn).DefaultScreen(conn).Root
	atomCookie := xproto.InternAtom(conn, true, uint16(len("CLIPBOARD")), "CLIPBOARD")
	atomReply, err := atomCookie.Reply()
	if err != nil {
		conn.Close()
		return nil, err
	}

	mask := uint32(xfixes.SelectionEventMaskSetSelectionOwner |
		xfixes.SelectionEventMaskSelectionWindowDestroy |
		xfixes.SelectionEventMaskSelectionClientClose)

	xfixes.SelectSelectionInput(conn, root, atomReply.Atom, mask)

	w := &X11EventWatcher{
		conn:   conn,
		events: make(chan struct{}, 1),
		errs:   make(chan error, 1),
	}

	go w.loop()
	return w, nil
}

// Events returns a channel that receives on clipboard changes.
func (w *X11EventWatcher) Events() <-chan struct{} {
	return w.events
}

// Errors returns a channel that receives async errors.
func (w *X11EventWatcher) Errors() <-chan error {
	return w.errs
}

// Close releases the X11 connection.
func (w *X11EventWatcher) Close() error {
	if w.conn == nil {
		return nil
	}
	w.conn.Close()
	w.conn = nil
	return nil
}

func (w *X11EventWatcher) loop() {
	for {
		event, err := w.conn.WaitForEvent()
		if err != nil {
			select {
			case w.errs <- err:
			default:
			}
			close(w.events)
			return
		}
		if event == nil {
			close(w.events)
			return
		}

		switch event.(type) {
		case xfixes.SelectionNotifyEvent:
			select {
			case w.events <- struct{}{}:
			default:
			}
		}
	}
}
