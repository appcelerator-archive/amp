package user

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type signUpOpts struct {
	username string
	email    string
	password string
}

var (
	opts = &signUpOpts{}
)

// NewSignUpCommand returns a new instance of the signup command.
func NewSignUpCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "signup [OPTIONS]",
		Short:   "Signup for a new account",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return signUp(c, cmd)
		},
	}
	cmd.Flags().StringVar(&opts.username, "name", "", "User name")
	cmd.Flags().StringVar(&opts.email, "email", "", "User email")
	cmd.Flags().StringVar(&opts.password, "password", "", "User password")
	return cmd
}

func signUp(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("name").Changed {
		opts.username = c.Console().GetInput("username")
	}
	if !cmd.Flag("email").Changed {
		opts.email = c.Console().GetInput("email")
	}
	if !cmd.Flag("password").Changed {
		opts.password = c.Console().GetSilentInput("password")
	}

	conn, err := c.ClientConn()
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.SignUpRequest{
		Name:     opts.username,
		Email:    opts.email,
		Password: opts.password,
	}
	_, err = client.SignUp(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Printf("Hi %s! Please check your email to complete the signup process.\n", opts.username)
	return nil
}
