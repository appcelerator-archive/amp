package main

import (
	"github.com/spf13/cobra"
)

// ServiceCmd is the main command for attaching service subcommands.
var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage services",
	Long:  `Manage services.`,
}

func init() {
	RootCmd.AddCommand(ServiceCmd)
}
