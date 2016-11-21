package main

import (
	"github.com/spf13/cobra"
	"os"
)

// PlatformStart is the main command for attaching platform subcommands.
var PlatformStart = &cobra.Command{
	Use:   "start",
	Short: "Start platform",
	Long:  `Start all AMP platform services.`,
	Run: func(cmd *cobra.Command, args []string) {
		startAMP(cmd, args)
	},
}

func init() {
	PlatformStart.Flags().BoolP("force", "f", false, "Start all possible services, do not stop on error")
	PlatformStart.Flags().BoolP("silence", "s", false, "no console output at all")
	PlatformCmd.AddCommand(PlatformStart)
}

func startAMP(cmd *cobra.Command, args []string) error {
	manager := &ampManager{}
	if cmd.Flag("silence").Value.String() == "true" {
		manager.silence = true
	}
	if cmd.Flag("force").Value.String() == "true" {
		manager.force = true
	}
	if cmd.Flag("verbose").Value.String() == "true" {
		manager.verbose = true
	}
	if err := manager.init("Starting AMP platform"); err != nil {
		manager.printf(colError, "Start error: %v\n", err)
		os.Exit(1)
	}
	if cmd.Flag("server").Value.String() != "" {
		manager.printf(colWarn, "Error: --server has no effect for start command\n")
		os.Exit(1)
	}
	stack := getAMPInfrastructureStack(manager)
	manager.computeStatus(stack)
	if manager.status == "running" {
		if !manager.force {
			manager.printf(colMagenta, "AMP platform already started (-f to force a re-start)\n")
			return nil
		}
		if err := manager.stop(stack); err != nil {
			manager.printf(colWarn, "Stop error: %v\n", err)
			manager.printf(colWarn, "Mode force: start anyway\n")
		}
	}
	if err := manager.start(stack); err != nil {
		manager.printf(colError, "Start error: %v\n", err)
		if err := manager.stop(stack); err != nil {
			manager.printf(colError, "Stop error: %v\n", err)
		}
		os.Exit(1)
	}
	manager.printf(colMagenta, "AMP platform started\n")
	return nil
}
