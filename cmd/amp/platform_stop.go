package main

import (
	"os"

	"github.com/spf13/cobra"
)

// PlatformStop is the main command for attaching platform subcommands.
var PlatformStop = &cobra.Command{
	Use:   "stop",
	Short: "Stop platform",
	Long:  `Stop all AMP platform services.`,
	Run: func(cmd *cobra.Command, args []string) {
		stopAMP(cmd, args)
	},
}

func init() {
	PlatformStop.Flags().BoolP("quiet", "q", false, "Suppress terminal output")
	PlatformStop.Flags().BoolP("local", "l", false, "Use local amp image")
	PlatformCmd.AddCommand(PlatformStop)
}

func stopAMP(cmd *cobra.Command, args []string) error {
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
	if err := manager.init("Stopping AMP platform"); err != nil {
		manager.printf(colError, "Start error: %v\n", err)
		os.Exit(1)
	}
	if cmd.Flag("server").Value.String() != "" {
		manager.printf(colWarn, "Error: --server has no effect for stop command\n")
		os.Exit(1)
	}
	stack := getAMPInfrastructureStack(manager)
	manager.computeStatus(stack)
	if manager.status == "stopped" {
		manager.printf(colRegular, "AMP platform already stopped\n")
		return nil
	}
	if err := manager.stop(stack); err != nil {
		manager.printf(colError, "Stop error: %v\n", err)
		os.Exit(1)
	}
	manager.printf(colRegular, "AMP platform stopped\n")
	return nil
}
