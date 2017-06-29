package login

import (
	"errors"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type loginOptions struct {
	username string
	password string
}

// NewLoginCommand returns a new instance of the login command.
func NewLoginCommand(c cli.Interface) *cobra.Command {
	opts := loginOptions{}
	cmd := &cobra.Command{
		Use:     "login [OPTIONS]",
		Short:   "Login to account",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return login(c, cmd, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.username, "name", "", "User name")
	flags.StringVar(&opts.password, "password", "", "User password")
	return cmd
}

func login(c cli.Interface, cmd *cobra.Command, opts loginOptions) error {
	if !cmd.Flag("name").Changed {
		opts.username = c.Console().GetInput("username")
	}
	if !cmd.Flag("password").Changed {
		opts.password = c.Console().GetSilentInput("password")
	}

	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.LogInRequest{
		Name:     opts.username,
		Password: opts.password,
	}
	headers := metadata.MD{}
	_, err := client.Login(context.Background(), request, grpc.Header(&headers))
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	if err := cli.SaveToken(headers, c.Server()); err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	c.Console().Printf("Welcome back %s!\n", opts.username)
	return nil
}
