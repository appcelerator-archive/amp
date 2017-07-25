package secret

import (
	"github.com/spf13/cobra"
	"github.com/appcelerator/amp/cli"
)

// NewSecretCommand returns a new instance of the secret command.
func NewSecretCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "secret",
		Short:   "Manage secrets",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}

	cmd.AddCommand(NewCreateCommand(c))
	cmd.AddCommand(NewListCommand(c))
	cmd.AddCommand(NewRemoveCommand(c))

	return cmd
}

