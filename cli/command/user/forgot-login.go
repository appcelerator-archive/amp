package user

import (
	"errors"
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// NewForgotLoginCommand returns a new instance of the forgot-login command.
func NewForgotLoginCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "forgot-login EMAIL",
		Short:   "Retrieve account name",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return forgotLogin(c, args)
		},
	}
}

func forgotLogin(c cli.Interface, args []string) error {
	if token := cli.GetToken(); token != "" {
		return errors.New("you are already logged into an account. Use 'amp whoami' to view your username")
	}
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.ForgotLoginRequest{
		Email: args[0],
	}
	if _, err := client.ForgotLogin(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Printf("Your login name has been sent to the address: %s\n", args[0])
	return nil
}
