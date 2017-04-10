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

type verifyOpts struct {
	token string
}

var (
	verifyOptions = &verifyOpts{}
)

// NewVerifyCommand returns a new instance of the verify command.
func NewVerifyCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "verify TOKEN",
		Short:   "Verify account",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return errors.New("token cannot be empty")
			}
			verifyOptions.token = args[0]
			return verify(c, verifyOptions)
		},
	}
}

func verify(c cli.Interface, opt *verifyOpts) error {
	conn, err := c.ClientConn()
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.VerificationRequest{
		Token: opt.token,
	}
	_, err = client.Verify(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Your account has now been activated.")
	return nil
}
