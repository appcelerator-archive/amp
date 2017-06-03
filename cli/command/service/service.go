package service

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewServiceCommand returns a new instance of the service command.
func NewServiceCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service",
		Short:   "Service management operations",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}
	cmd.AddCommand(NewServiceListCommand(c))
	cmd.AddCommand(NewServiceLogsCommand(c))
	cmd.AddCommand(NewServicePsCommand(c))
	return cmd
}
