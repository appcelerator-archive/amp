package password

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type changeOpts struct {
	current_password string
	new_password     string
}

var (
	changeOptions = &changeOpts{}
)

// NewChangeCommand returns a new instance of the change command.
func NewChangeCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "change [OPTIONS]",
		Short:   "Change password",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return change(c, cmd)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&changeOptions.current_password, "current", "", "Current password")
	flags.StringVar(&changeOptions.new_password, "new", "", "New password")
	return cmd
}

func change(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("current").Changed {
		changeOptions.current_password = c.Console().GetSilentInput("current password")
	}
	if !cmd.Flag("new").Changed {
		changeOptions.new_password = c.Console().GetSilentInput("new password")
	}

	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.PasswordChangeRequest{
		ExistingPassword: changeOptions.current_password,
		NewPassword:      changeOptions.new_password,
	}
	if _, err := client.PasswordChange(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Your password change has been successful.")
	return nil
}
