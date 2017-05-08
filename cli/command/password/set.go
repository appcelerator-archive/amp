package password

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type setPasswordOptions struct {
	token    string
	password string
}

// NewSetCommand returns a new instance of the set command.
func NewSetCommand(c cli.Interface) *cobra.Command {
	opts := setPasswordOptions{}
	cmd := &cobra.Command{
		Use:     "set [OPTIONS]",
		Short:   "Set password",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return set(c, cmd, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.token, "token", "", "Verification token")
	flags.StringVar(&opts.password, "password", "", "User password")
	return cmd
}

func set(c cli.Interface, cmd *cobra.Command, opts setPasswordOptions) error {
	if !cmd.Flag("token").Changed {
		opts.token = c.Console().GetInput("token")
	}
	if !cmd.Flag("password").Changed {
		opts.password = c.Console().GetSilentInput("password")
	}

	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.PasswordSetRequest{
		Token:    opts.token,
		Password: opts.password,
	}
	if _, err := client.PasswordSet(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Your password set has been successful.")
	return nil
}
