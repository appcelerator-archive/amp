package cli

import (
	"testing"
	"time"

	"golang.org/x/net/context"
)

// createTimeout creates the timeout goroutine and returns the func to cancel it.
func createTimeout(t *testing.T, duration time.Duration, name string) (cancel func()) {
	// Create the context with the specified duration.
	ctx, cancel := context.WithTimeout(context.Background(), duration)

	// Create the checkTimeout goroutine.
	go checkTimeout(t, ctx, name)

	// return the cancel function to end the goroutine.
	return cancel
}

// checkTimeout repeatedly checks the timeout and returns or fails when the deadline is exceeded.
func checkTimeout(t *testing.T, ctx context.Context, name string) {
	// Loop.
	for {
		select {
		case <-ctx.Done():
			// If the deadline exceeds, fail the test.
			if ctx.Err() == context.DeadlineExceeded {
				t.Fatal("Deadline exceeded:", name)
			}
			// If the goroutine is cancelled, return.
			return
		}
	}
}
