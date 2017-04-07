package user

import (
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
	cmd := &cobra.Command{
		Use:     "verify",
		Short:   "Verify account",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return verify(c, cmd)
		},
	}
	cmd.Flags().StringVar(&verifyOptions.token, "token", "", "Verification token")
	return cmd
}

func verify(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("token").Changed {
		verifyOptions.token = c.Console().GetInput("token")
	}
	conn, err := c.ClientConn()
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.VerificationRequest{
		Token: verifyOptions.token,
	}
	_, err = client.Verify(context.Background(), request)
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	c.Console().Println("Your account has now been activated.")
	return nil
}
