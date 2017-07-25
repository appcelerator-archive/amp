package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string
)

func main() {
	log.Printf("ampctl (version: %s, build: %s)\n", Version, Build)
	rootCmd := &cobra.Command{
		Use:   "ampctl",
		Short: "Run commands in target amp cluster",
	}

	rootCmd.AddCommand(NewChecksCommand())
	rootCmd.AddCommand(NewInstallCommand())
	rootCmd.AddCommand(NewMonitorCommand())

	// if no arg has been provided, execute check and install
	if len(os.Args[1:]) == 0 {
		rootCmd.SetArgs([]string{"check"})
		err := rootCmd.Execute()
		if err != nil {
			log.Println(err)
			os.Exit(-1)
		}
		rootCmd.SetArgs([]string{"install"})
		err = rootCmd.Execute()
		if err != nil {
			log.Println(err)
			os.Exit(-1)
		}
	} else {
		err := rootCmd.Execute()
		if err != nil {
			log.Println(err)
			os.Exit(-1)
		}
	}
}
