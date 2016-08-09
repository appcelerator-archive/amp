package data

import "time"

// Store must be implemented for
type Stats interface {
	// Endpoints returns an array of endpoints for the storage
	Endpoints() []string

	// Connect to stats server
	Connect(timeout time.Duration) error

	// Close connection to stats server
	Close() error

	Query(query string) (string, error)
}
