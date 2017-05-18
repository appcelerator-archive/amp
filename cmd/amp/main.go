package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/appcelerator/amp/cli"
	"github.com/docker/docker/pkg/term"
)

var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string

	//address string
	config *cli.Configuration
)

func main() {
	// Initial configuration (before processing environment and flags)
	config = &cli.Configuration{
		Version: Version,
		Build:   Build,
		Server:  cli.DefaultAddress + cli.DefaultPort,
	}

	// Read the cli config
	if err := cli.ReadClientConfig(config); err != nil {
		log.Fatalln(err)
	}

	// Set terminal emulation based on platform as required.
	stdin, stdout, stderr := term.StdStreams()
	c := cli.NewCLI(stdin, stdout, stderr, config)
	cmd := newRootCommand(c)

	if err := cmd.Execute(); err != nil {
		// cobra command error looks like: unknown command "foo" for "amp"
		// replace it with error that is more consistent with errors for flags and args.
		e := err.Error()
		pattern := "unknown command \"([^\"]*)\""
		re := regexp.MustCompile(pattern)
		match := re.FindStringSubmatch(e)
		if match != nil {
			e = fmt.Sprintf("not an amp command: %s\nSee '%s --help", match[1], cmd.CommandPath())
		}

		c.Console().Errorln(e)
		cli.Exit(1)
	}
}

func showVersion() {
	fmt.Printf("amp version %s, build %s\n", Version, Build)
}
