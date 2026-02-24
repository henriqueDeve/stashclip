package clipboard

// ClipboardProvider provides clipboard read/write operations.
type ClipboardProvider interface {
	Read() (string, error)
	Write(text string) error
}
