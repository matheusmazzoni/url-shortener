package storage

import (
	"database/sql"

	// The blank import is used for the side-effect of registering the driver.
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStore is the implementation of the Storage interface for a SQLite database.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a new instance of SQLiteStore, opens a database connection,
// and ensures the necessary tables and indexes exist.
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	// sql.Open() doesn't immediately open a connection. It prepares a connection pool object.
	// The first actual connection is established lazily, when it's first needed.
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL") // WAL mode is better for concurrency.
	if err != nil {
		return nil, err
	}

	// Ping the database to verify the connection is alive and working.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS urls (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"short_key" TEXT NOT NULL UNIQUE,
		"original_url" TEXT NOT NULL
	);
	
	-- An index on original_url is crucial for the performance of GetShortKey lookups.
	CREATE INDEX IF NOT EXISTS idx_original_url ON urls(original_url);
	`

	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

// SaveURL saves a new URL mapping (short key to original URL) into the database.
func (s *SQLiteStore) SaveURL(shortKey, originalURL string) error {
	_, err := s.db.Exec("INSERT INTO urls (short_key, original_url) VALUES (?, ?)", shortKey, originalURL)
	return err
}

// GetURL retrieves an original URL from the database given a short key.
func (s *SQLiteStore) GetURL(shortKey string) (string, error) {
	var originalURL string
	err := s.db.QueryRow("SELECT original_url FROM urls WHERE short_key = ?", shortKey).Scan(&originalURL)
	return originalURL, err
}

// GetShortKey retrieves an existing short key from the database given an original URL.
func (s *SQLiteStore) GetShortKey(originalURL string) (string, error) {
	var shortKey string
	err := s.db.QueryRow("SELECT short_key FROM urls WHERE original_url = ?", originalURL).Scan(&shortKey)
	return shortKey, err
}

// KeyExists checks if a given short key already exists in the database.
func (s *SQLiteStore) KeyExists(shortKey string) (bool, error) {
	var exists int
	// We query for a constant '1' to check for the row's existence efficiently.
	// If the row doesn't exist, we'll get sql.ErrNoRows.
	err := s.db.QueryRow("SELECT 1 FROM urls WHERE short_key = ? LIMIT 1", shortKey).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			// This is not an application error. It means the key is available.
			return false, nil
		}
		// Any other error indicates a real problem with the database.
		return false, err
	}

	// If there was no error, the row was found, and the key exists.
	return true, nil
}

// Close closes the database connection, releasing all its resources.
// It's important to call this to prevent connection leaks.
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
