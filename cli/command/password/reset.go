package password

import (
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type resetOpts struct {
	username string
}

var (
	resetOptions = &resetOpts{}
)

// NewResetCommand returns a new instance of the reset command.
func NewResetCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "reset",
		Short:   "Reset password",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return reset(c, cmd)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&resetOptions.username, "name", "", "User name")
	return cmd
}

func reset(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("name").Changed {
		resetOptions.username = c.Console().GetInput("username")
	}

	conn, err := c.ClientConn()
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.PasswordResetRequest{
		Name: resetOptions.username,
	}
	if _, err = client.PasswordReset(context.Background(), request); err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	c.Console().Printf("Hi %s! Please check your email to complete the password reset process.\n", resetOptions.username)
	return nil
}
