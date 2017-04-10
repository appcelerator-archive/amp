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

type forgotOpts struct {
	email string
}

var (
	forgotOptions = &forgotOpts{}
)

// NewForgotLoginCommand returns a new instance of the forgot-login command.
func NewForgotLoginCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "forgot-login EMAIL",
		Short:   "Retrieve account name",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return errors.New("email cannot be empty")
			}
			forgotOptions.email = args[0]
			return forgotLogin(c, forgotOptions)
		},
	}
}

func forgotLogin(c cli.Interface, opt *forgotOpts) error {
	conn, err := c.ClientConn()
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.ForgotLoginRequest{
		Email: opt.email,
	}
	if _, err = client.ForgotLogin(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Printf("Your login name has been sent to the address: %s", opt.email)
	return nil
}
