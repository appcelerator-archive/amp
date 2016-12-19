package main

import (
	"os"

	"github.com/spf13/cobra"
)

// PlatformStatus is the main command for attaching platform subcommands.
var PlatformStatus = &cobra.Command{
	Use:   "status",
	Short: "Get AMP platform status",
	Long:  `Get AMP platform global status (stopped, partially running, running command return 1 if status is not running.`,
	Run: func(cmd *cobra.Command, args []string) {
		getAMPStatus(cmd, args)
	},
}

func init() {
	PlatformStatus.Flags().BoolP("quiet", "q", false, "Suppress terminal output")
	PlatformStatus.Flags().BoolP("local", "l", false, "Use local amp image")
	PlatformCmd.AddCommand(PlatformStatus)
}

func getAMPStatus(cmd *cobra.Command, args []string) error {
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
	if err := manager.init(""); err != nil {
		manager.printf(colError, "Compute status error: %v\n", err)
		os.Exit(1)
	}
	if cmd.Flag("server").Value.String() != "" {
		manager.printf(colWarn, "Error: --server has no effect for status command\n")
		os.Exit(1)
	}
	manager.computeStatus(getAMPInfrastructureStack(manager))
	manager.printf(colRegular, "status: %s\n", manager.status)
	if manager.status != "running" {
		os.Exit(1)
	}
	return nil
}
