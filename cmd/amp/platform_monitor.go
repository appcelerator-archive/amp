package main

import (
	"github.com/spf13/cobra"
	"os"
)

// PlatformMonitor is the main command for attaching platform subcommands.
var PlatformMonitor = &cobra.Command{
	Use:   "monitor",
	Short: "Display AMP platform services",
	Long:  `Display AMP platform services information and states.`,
	Run: func(cmd *cobra.Command, args []string) {
		displayAMPServiceStatus(cmd, args)
	},
}

func init() {
	PlatformCmd.AddCommand(PlatformMonitor)
}

func displayAMPServiceStatus(cmd *cobra.Command, args []string) error {
	manager := &ampManager{}
	if cmd.Flag("verbose").Value.String() == "true" {
		manager.verbose = true
	}
	if err := manager.init(""); err != nil {
		manager.printf(colError, "Monitor error: %v\n", err)
		os.Exit(1)
	}
	if cmd.Flag("server").Value.String() != "" {
		manager.printf(colWarn, "Error: --server has no effect for monitor command\n")
		os.Exit(1)
	}
	manager.monitor(getAMPInfrastructureStack(manager))
	return nil
}
