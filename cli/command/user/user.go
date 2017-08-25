package user

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewUserCommand returns a new instance of the user command.
func NewUserCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user",
		Short:   "User management operations",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}
	cmd.AddCommand(NewSignUpCommand(c))
	cmd.AddCommand(NewVerifyCommand(c))
	cmd.AddCommand(NewResendVerifyCommand(c))
	cmd.AddCommand(NewForgotLoginCommand(c))
	cmd.AddCommand(NewListUserCommand(c))
	cmd.AddCommand(NewGetUserCommand(c))
	cmd.AddCommand(NewRemoveUserCommand(c))
	return cmd
}
