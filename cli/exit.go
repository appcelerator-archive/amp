package cli

import (
	"os"
)

var (
	onExitFuncs []func()
)

// OnExit is used to register functions that should be executed upon exiting.
func OnExit(f func()) {
	onExitFuncs = append(onExitFuncs, f)
}

// Exit ensures that any registered functions are executed before exiting
// with the specified status code.
func Exit(code int) {
	for _, f := range onExitFuncs {
		f()
	}
	os.Exit(code)
}

