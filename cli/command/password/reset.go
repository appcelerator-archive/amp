package password

import (
	"errors"
	"fmt"

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
	return &cobra.Command{
		Use:     "reset USERNAME",
		Short:   "Reset password",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return errors.New("username cannot be empty")
			}
			resetOptions.username = args[0]
			return reset(c, resetOptions)
		},
	}
}

func reset(c cli.Interface, opt *resetOpts) error {
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.PasswordResetRequest{
		Name: opt.username,
	}
	if _, err := client.PasswordReset(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Printf("Hi %s! Please check your email to complete the password reset process.\n", opt.username)
	return nil
}
