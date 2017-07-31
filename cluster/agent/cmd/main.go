package main

import (
	"log"

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
		RunE: func(cmd *cobra.Command, args []string) error {
			// perform checks and install by default when no sub-command is specified
			if err := checks(cmd, []string{}); err != nil {
				return err
			}
			return install(cmd, args)
		},
	}

	rootCmd.AddCommand(NewChecksCommand())
	rootCmd.AddCommand(NewInstallCommand())
	rootCmd.AddCommand(NewMonitorCommand())
	rootCmd.AddCommand(NewUninstallCommand())

	// should be in the Install command, but since the local
	rootCmd.PersistentFlags().BoolVar(&installOpts.skipTests, "fast", false, "Skip tests while deploying the core services")
	rootCmd.PersistentFlags().BoolVar(&installOpts.noMonitoring, "no-monitoring", false, "Don't deploy the monitoring core services")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %s\n", err)
	}
}
