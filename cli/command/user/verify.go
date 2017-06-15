package user

import (
	"errors"
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/api/rpc/version"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
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
	clientVer := version.NewVersionClient(conn)
	requestVer := &version.GetRequest{}
	reply, err := clientVer.Get(context.Background(), requestVer)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if reply.Info.Registration == configuration.RegistrationNone {
		return errors.New("`amp user verify` disabled. This cluster has no registration policy")
	}

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
