package login

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type loginOpts struct {
	username string
	email    string
	password string
}

var (
	loginOptions = &loginOpts{}
)

// NewLoginCommand returns a new instance of the login command.
func NewLoginCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "login [OPTIONS]",
		Short:   "Login to account",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return login(c, cmd)
		},
	}
	cmd.Flags().StringVar(&loginOptions.username, "name", "", "User name")
	cmd.Flags().StringVar(&loginOptions.password, "password", "", "User password")
	return cmd
}

func login(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("name").Changed {
		loginOptions.username = c.Console().GetInput("username")
	}
	if !cmd.Flag("password").Changed {
		loginOptions.password = c.Console().GetSilentInput("password")
	}

	conn, err := c.ClientConn()
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.LogInRequest{
		Name:     loginOptions.username,
		Password: loginOptions.password,
	}
	header := metadata.MD{}
	_, err = client.Login(context.Background(), request, grpc.Header(&header))
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if err := cli.SaveToken(header); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Printf("Welcome back %s!\n", loginOptions.username)
	return nil
}
