package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	atexitFuncs []func()
)

// Exit ensures that any registered functions are executed before exiting
// with the specified status code.
func Exit(code int) {
	for _, f := range atexitFuncs {
		f()
	}
	os.Exit(code)
}

// AtExit is used to register functions to execute before exiting.
func AtExit(f func()) {
	atexitFuncs = append(atexitFuncs, f)
}

// PrintErr prints the error and then exits
func PrintErr(err error) {
	color.Set(color.FgRed)
	fmt.Println(err)
	os.Exit(1)
}
