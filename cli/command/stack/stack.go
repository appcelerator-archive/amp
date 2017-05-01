package stack

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewStackCommand returns a new instance of the stack command.
func NewStackCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stack",
		Short:   "Stack management operations",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}
	cmd.AddCommand(NewDeployCommand(c))
	cmd.AddCommand(NewListCommand(c))
	cmd.AddCommand(NewRemoveCommand(c))
	cmd.AddCommand(NewLogsCommand(c))
	cmd.AddCommand(NewServicesCommand(c))
	return cmd
}
