package function_

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewFunctionCommand returns a new instance of the function command.
func NewFunctionCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "function",
		Short:   "Function management operations",
		Aliases: []string{"fn"},
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}
	cmd.AddCommand(NewFunctionCreateCommand(c))
	cmd.AddCommand(NewFunctionListCommand(c))
	cmd.AddCommand(NewFunctionRemoveCommand(c))
	return cmd
}
