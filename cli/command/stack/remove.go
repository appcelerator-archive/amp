package stack

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewStackCommand returns a new instance of the stack command.
func NewRemoveCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Remove a deployed stack",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}
	return cmd
}
