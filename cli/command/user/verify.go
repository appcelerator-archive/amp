package user

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// NewVerifyCommand returns a new instance of the verify command.
func NewVerifyCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "verify TOKEN",
		Short:   "Verify account",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return verify(c, args)
		},
	}
}

func verify(c cli.Interface, args []string) error {
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.VerificationRequest{
		Token: args[0],
	}
	if _, err := client.Verify(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Your account has now been activated.")
	return nil
}
