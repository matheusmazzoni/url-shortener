package storage

// Storage defines the interface for URL storage operations.
// Any database type that wants to be used must implement this interface.
type Storage interface {
	SaveURL(shortKey, originalURL string) error
	GetURL(shortKey string) (string, error)
	GetShortKey(originalURL string) (string, error)
	KeyExists(shortKey string) (bool, error)
	Close() error
}
