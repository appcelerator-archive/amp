package config

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewConfigCommand returns a new instance of the config command.
func NewConfigCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Short:   "Configuration management operations",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}

	cmd.AddCommand(NewCreateCommand(c))
	cmd.AddCommand(NewListCommand(c))
	cmd.AddCommand(NewRemoveCommand(c))

	return cmd
}
