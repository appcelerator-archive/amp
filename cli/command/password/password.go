package password

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewPasswordCommand returns a new instance of the password command.
func NewPasswordCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "password",
		Short:   "Password management operations",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}
	cmd.AddCommand(NewChangeCommand(c))
	cmd.AddCommand(NewResetCommand(c))
	cmd.AddCommand(NewSetCommand(c))
	return cmd
}
