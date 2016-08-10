package stat

import "time"

// Stats must be implemented for
type Stats interface {

	// Connect to stats server
	Connect(timeout time.Duration) error

	// Close connection to stats server
	Close() error

	Query(query string) (string, error)
}
