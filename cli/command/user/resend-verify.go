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

// NewResendVerifyCommand returns a new instance of the resend-verify command.
func NewResendVerifyCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "resend-verification-token USERNAME",
		Short:   "Resend verification email to registered address",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resendVerify(c, args)
		},
	}
}

func resendVerify(c cli.Interface, args []string) error {
	conn := c.ClientConn()
	clientVer := version.NewVersionClient(conn)
	requestVer := &version.GetRequest{}
	reply, err := clientVer.VersionGet(context.Background(), requestVer)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	if reply.Info.Registration == configuration.RegistrationNone {
		return errors.New("`amp user resend-verify` disabled. This cluster has no registration policy")
	}

	client := account.NewAccountClient(conn)
	request := &account.ResendVerifyRequest{
		Name: args[0],
		Url:  c.Server(),
	}
	if _, err := client.ResendVerify(context.Background(), request); err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	c.Console().Printf("A new verification email has been sent to %s\n", args[0])
	return nil
}
