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

// PrintErr prints the error and then exits
func printErr(err error) {
	color.Set(color.FgRed)
	fmt.Println(err)
	os.Exit(1)
}

// SaveConfiguration saves the configuration to ~/.ampswarm.yaml
func saveConfiguration(c interface{}) (err error) {
	homedir, err := homedir.Dir()
	if err != nil {
		return
	}
	contents, err := yaml.Marshal(c)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path.Join(homedir, ".amp.yaml"), contents, os.ModePerm)
	if err != nil {
		return
	}
	return
}
