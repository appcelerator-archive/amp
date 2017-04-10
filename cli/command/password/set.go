package password

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type setOpts struct {
	token    string
	password string
}

var (
	setOptions = &setOpts{}
)

// NewSetCommand returns a new instance of the set command.
func NewSetCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set [OPTIONS]",
		Short:   "Set password",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return set(c, cmd)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&setOptions.token, "token", "", "Verification token")
	flags.StringVar(&setOptions.password, "password", "", "User password")
	return cmd
}

func set(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("token").Changed {
		setOptions.token = c.Console().GetInput("token")
	}
	if !cmd.Flag("password").Changed {
		setOptions.password = c.Console().GetSilentInput("password")
	}

	conn, err := c.ClientConn()
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.PasswordSetRequest{
		Token:    setOptions.token,
		Password: setOptions.password,
	}
	if _, err = client.PasswordSet(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Your password set has been successful.")
	return nil
}
