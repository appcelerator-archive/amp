package password

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// NewResetCommand returns a new instance of the reset command.
func NewResetCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "reset USERNAME",
		Short:   "Reset password",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return reset(c, args)
		},
	}
}

func reset(c cli.Interface, args []string) error {
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.PasswordResetRequest{
		Name: args[0],
	}
	if _, err := client.PasswordReset(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Printf("Hi %s! Please check your email to complete the password reset process.\n", args[0])
	return nil
}
