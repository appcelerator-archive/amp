package main

import (
	"log"
	"os"

	"github.com/appcelerator/amp/cluster/ampagent/cmd"
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
		RunE: func(command *cobra.Command, args []string) error {
			// perform checks and install by default when no sub-command is specified
			if err := cmd.Checks(command, []string{}); err != nil {
				return err
			}
			return cmd.Install(command, args)
		},
	}

	rootCmd.AddCommand(cmd.NewChecksCommand())
	rootCmd.AddCommand(cmd.NewInstallCommand())
	rootCmd.AddCommand(cmd.NewUninstallCommand())

	if os.Getenv("TAG") == "" { // If TAG is undefined, use the current project version
		os.Setenv("TAG", Version)
	}
	if val, ok := os.LookupEnv("NO_LOGS"); ok && val == "true" {
		cmd.InstallOpts.NoLogs = true
	}
	if val, ok := os.LookupEnv("NO_METRICS"); ok && val == "true" {
		cmd.InstallOpts.NoMetrics = true
	}
	if val, ok := os.LookupEnv("NO_PROXY"); ok && val == "true" {
		cmd.InstallOpts.NoProxy = true
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %s\n", err)
	}
}
