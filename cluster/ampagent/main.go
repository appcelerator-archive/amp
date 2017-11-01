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

	// These flags pertain to install, but need to be enabled here at root and persist for when it is invoked with no subcommand
	rootCmd.PersistentFlags().BoolVar(&cmd.InstallOpts.NoLogs, "no-logs", false, "Don't deploy logs stack")
	rootCmd.PersistentFlags().BoolVar(&cmd.InstallOpts.NoMetrics, "no-metrics", false, "Don't deploy metrics stack")
	rootCmd.PersistentFlags().BoolVar(&cmd.InstallOpts.NoProxy, "no-proxy", false, "Don't deploy proxy stack")
	rootCmd.PersistentFlags().BoolVar(&cmd.InstallOpts.NoNodeManagement, "no-node-management", false, "Don't deploy node management stack")

	// Environment variables
	if os.Getenv("TAG") == "" { // If TAG is undefined, use the current project version
		os.Setenv("TAG", Version)
	}
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %s\n", err)
	}
}
