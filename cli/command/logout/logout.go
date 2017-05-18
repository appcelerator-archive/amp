package logout

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewLogoutCommand returns a new instance of the logout command.
func NewLogoutCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "logout",
		Short:   "Logout of account",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return logout(c)
		},
	}
}

func logout(c cli.Interface) error {
	if err := cli.RemoveToken(c.Server()); err != nil {
		return err
	}
	c.Console().Println("You have been logged out!")
	return nil
}
