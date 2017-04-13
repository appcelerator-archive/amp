package stack

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewStackCommand returns a new instance of the stack command.
func NewListCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List deployed stacks",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}
	return cmd
}
