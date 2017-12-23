package main

import (
	"log"
	"os"

	"github.com/appcelerator/amp/cluster/ampagent/cmd"
)

var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string

	// Default pattern for the overlay network creation
	DefaultSubnetPattern = "10.%d.0.0/16"
)

func main() {
	log.Printf("ampctl (version: %s, build: %s)\n", Version, Build)
	rootCmd := cmd.NewRootCommand()
	rootCmd.AddCommand(cmd.NewChecksCommand())
	rootCmd.AddCommand(cmd.NewInstallCommand())
	rootCmd.AddCommand(cmd.NewUninstallCommand())

	// These flags pertain to install, but need to be enabled here at root and persist for when it is invoked with no subcommand
	rootCmd.PersistentFlags().BoolVar(&cmd.InstallOpts.NoLogs, "no-logs", false, "Don't deploy logs stack")
	rootCmd.PersistentFlags().BoolVar(&cmd.InstallOpts.NoMetrics, "no-metrics", false, "Don't deploy metrics stack")
	rootCmd.PersistentFlags().BoolVar(&cmd.InstallOpts.NoProxy, "no-proxy", false, "Don't deploy proxy stack")
	rootCmd.PersistentFlags().StringVar(&cmd.InstallOpts.SubnetPattern, "subnet-pattern", DefaultSubnetPattern, "Subnet pattern for overlay networks, should contain a single %d")

	// Environment variables
	if os.Getenv("TAG") == "" { // If TAG is undefined, use the current project version
		os.Setenv("TAG", Version)
	}
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %s\n", err)
	}
}
