package password

import (
	"errors"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

type changePasswordOptions struct {
	current_password string
	new_password     string
}

// NewChangeCommand returns a new instance of the change command.
func NewChangeCommand(c cli.Interface) *cobra.Command {
	opts := changePasswordOptions{}
	cmd := &cobra.Command{
		Use:     "change [OPTIONS]",
		Short:   "Change password",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return change(c, cmd, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.current_password, "current", "", "Current password")
	flags.StringVar(&opts.new_password, "new", "", "New password")
	return cmd
}

func change(c cli.Interface, cmd *cobra.Command, opts changePasswordOptions) error {
	if !cmd.Flag("current").Changed {
		opts.current_password = c.Console().GetSilentInput("current password")
	}
	if !cmd.Flag("new").Changed {
		opts.new_password = c.Console().GetSilentInput("new password")
	}

	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.PasswordChangeRequest{
		ExistingPassword: opts.current_password,
		NewPassword:      opts.new_password,
	}
	if _, err := client.PasswordChange(context.Background(), request); err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	c.Console().Println("Your password change has been successful.")
	return nil
}
