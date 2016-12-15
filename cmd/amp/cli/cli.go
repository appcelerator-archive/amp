package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
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

// SaveConfiguration saves the configuration to ~/.config/amp/amp.yaml
func SaveConfiguration(c interface{}) (err error) {
	var configdir string
	xdgdir := os.Getenv("XDG_CONFIG_HOME")
	if xdgdir != "" {
		configdir = path.Join(xdgdir, "amp")
	} else {
		homedir, err := homedir.Dir()
		if err != nil {
			return err
		}
		configdir = path.Join(homedir, ".config/amp")
	}
	err = os.MkdirAll(configdir, 0755)
	if err != nil {
		return
	}
	contents, err := yaml.Marshal(c)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path.Join(configdir, "amp.yaml"), contents, os.ModePerm)
	if err != nil {
		return
	}
	return
}
