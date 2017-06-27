package user

import (
	"errors"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/api/rpc/version"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
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
	clientVer := version.NewVersionClient(conn)
	requestVer := &version.GetRequest{}
	reply, err := clientVer.Get(context.Background(), requestVer)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	if reply.Info.Registration == configuration.RegistrationNone {
		return errors.New("`amp user verify` disabled. This cluster has no registration policy")
	}

	client := account.NewAccountClient(conn)
	request := &account.VerificationRequest{
		Token: args[0],
	}
	if _, err := client.Verify(context.Background(), request); err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	c.Console().Println("Your account has now been activated.")
	return nil
}
