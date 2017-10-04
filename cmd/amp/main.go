package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/docker/docker/pkg/term"
)

var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string

	//address string
	cfg *cli.Configuration
)

func main() {
	// Initial configuration (before processing environment and flags)
	cfg = &cli.Configuration{
		Server: cli.DefaultAddress + cli.DefaultPort,
	}

	// Read the cli config
	if err := cli.ReadClientConfig(cfg); err != nil {
		log.Fatalln(err)
	}

	// Build and Version cannot be overridden
	cfg.Build = Build
	cfg.Version = Version

	// Set terminal emulation based on platform as required.
	stdin, stdout, stderr := term.StdStreams()
	c := cli.NewCLI(stdin, stdout, stderr, cfg)
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
