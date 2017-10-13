// NOTICE
// This file is copyrighted by Docker under the Apache 2.0 license
// The original can be found here:
// https://github.com/appcelerator/go-docker/docker/blob/master/cli/error.go

package cli

import (
	"strings"
)

// Errors is a list of errors.
// Useful in a loop if you don't want to return the error right away and you want to display after the loop,
// all the errors that happened during the loop.
type Errors []error

func (errList Errors) Error() string {
	if len(errList) < 1 {
		return ""
	}

	out := make([]string, len(errList))
	for i := range errList {
		out[i] = errList[i].Error()
	}
	return strings.Join(out, ", ")
}

// StatusError reports an unsuccessful exit by a command.
type StatusError struct {
	Status     string
	StatusCode int
}

func (e StatusError) Error() string {
	// TODO revisit the idea of command status codes, not useful right now
	// return fmt.Sprintf("Status: %s, Code: %d", e.Status, e.StatusCode)
	return e.Status
}
