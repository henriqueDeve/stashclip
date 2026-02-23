package clipboard

// Backend provides clipboard read/write operations.
type Backend interface {
	Read() (string, error)
	Write(text string) error
}
