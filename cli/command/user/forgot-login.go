package user

import (
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type forgotLoginOpts struct {
	email string
}

var (
	forgotLoginOptions = &forgotLoginOpts{}
)

// NewForgotLoginCommand returns a new instance of the forgot-login command.
func NewForgotLoginCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "forgot-login",
		Short:   "Retrieve account name",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return forgotLogin(c, cmd)
		},
	}

	cmd.Flags().StringVar(&forgotLoginOptions.email, "email", "", "User email")
	return cmd
}

func forgotLogin(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("email").Changed {
		forgotLoginOptions.email = c.Console().GetInput("email")
	}

	conn, err := c.ClientConn()
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.ForgotLoginRequest{
		Email: forgotLoginOptions.email,
	}
	if _, err = client.ForgotLogin(context.Background(), request); err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	c.Console().Printf("Your login name has been sent to the address: %s", forgotLoginOptions.email)
	return nil
}
