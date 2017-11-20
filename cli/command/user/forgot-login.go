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
		return errors.New("`amp user forgot-login` disabled. This cluster has no registration policy")
	}

	if token := cli.GetToken(c.Server()); token != "" {
		return errors.New("you are already logged into an account. Use 'amp whoami' to view your username")
	}
	client := account.NewAccountClient(conn)
	request := &account.ForgotLoginRequest{
		Email: args[0],
	}
	if _, err := client.ForgotLogin(context.Background(), request); err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	c.Console().Printf("Your login name has been sent to the address: %s\n", args[0])
	return nil
}
