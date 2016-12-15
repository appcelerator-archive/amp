package main

import (
	"os"

	"github.com/spf13/cobra"
)

// PlatformPull is the main command for attaching platform subcommands.
var PlatformPull = &cobra.Command{
	Use:   "pull",
	Short: "Pull platform images",
	Long:  `Pull all AMP platform images.`,
	Run: func(cmd *cobra.Command, args []string) {
		pullAMPImages(cmd, args)
	},
}

func init() {
	PlatformPull.Flags().BoolP("quiet", "q", false, "Suppress terminal output")
	PlatformPull.Flags().BoolP("local", "l", false, "Use local amp image")
	PlatformCmd.AddCommand(PlatformPull)
}

func pullAMPImages(cmd *cobra.Command, args []string) error {
	manager := &ampManager{}
	if cmd.Flag("quiet").Value.String() == "true" {
		manager.silence = true
	}
	if cmd.Flag("verbose").Value.String() == "true" {
		manager.verbose = true
	}
	if cmd.Flag("local").Value.String() == "true" {
		manager.local = true
	}
	if err := manager.init("Pulling AMP images"); err != nil {
		manager.printf(colError, "Pull error: %v\n", err)
		os.Exit(1)
	}
	if cmd.Flag("server").Value.String() != "" {
		manager.printf(colWarn, "Error: --server has no effect for pull command\n")
		os.Exit(1)
	}
	if err := manager.pull(getAMPInfrastructureStack(manager)); err != nil {
		manager.printf(colError, "Pull error: %v\n", err)
		os.Exit(1)
	}
	manager.printf(colRegular, "AMP platform images pulled\n")
	return nil
}
