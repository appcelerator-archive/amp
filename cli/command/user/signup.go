package user

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/api/rpc/version"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type signUpOptions struct {
	username  string
	email     string
	password  string
	autologin bool
}

// NewSignUpCommand returns a new instance of the signup command.
func NewSignUpCommand(c cli.Interface) *cobra.Command {
	opts := signUpOptions{}
	cmd := &cobra.Command{
		Use:     "signup [OPTIONS]",
		Short:   "Signup for a new account",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return signUp(c, cmd, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.username, "name", "", "User name")
	flags.StringVar(&opts.email, "email", "", "User email")
	flags.StringVar(&opts.password, "password", "", "User password")
	flags.BoolVar(&opts.autologin, "autologin", true, "Auto login to account, for clusters with no registration")
	return cmd
}

func signUp(c cli.Interface, cmd *cobra.Command, opts signUpOptions) error {
	if !cmd.Flag("name").Changed {
		opts.username = c.Console().GetInput("username")
	}
	if !cmd.Flag("email").Changed {
		opts.email = c.Console().GetInput("email")
	}
	if !cmd.Flag("password").Changed {
		opts.password = c.Console().GetSilentInput("password")
	}

	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	requestSignUp := &account.SignUpRequest{
		Name:     opts.username,
		Email:    opts.email,
		Password: opts.password,
	}
	if _, err := client.SignUp(context.Background(), requestSignUp); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}

	clientVer := version.NewVersionClient(conn)
	requestVer := &version.GetRequest{}
	reply, err := clientVer.Get(context.Background(), requestVer)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}

	// Auto login for cluster with no registration
	if reply.Info.Registration == "email" {
		c.Console().Printf("Hi %s! Please check your email to complete the signup process.\n", opts.username)
		return nil
	}
	c.Console().Printf("Verification is not necessary for this cluster.\n")
	if !opts.autologin {
		c.Console().Printf("Hi %s! Please log in with your details using `amp login`.\n", opts.username)
		return nil
	}

	requestLogin := &account.LogInRequest{
		Name:     opts.username,
		Password: opts.password,
	}
	headers := metadata.MD{}
	_, err = client.Login(context.Background(), requestLogin, grpc.Header(&headers))
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if err := cli.SaveToken(headers, c.Server()); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Printf("Hi %s! You have been automatically logged in.\n", opts.username)
	return nil
}
