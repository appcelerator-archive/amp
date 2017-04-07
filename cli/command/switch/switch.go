package switch_

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type switchOpts struct {
	account string
}

var (
	switchOptions = &switchOpts{}
)

// NewSwitchCommand returns a new instance of the switch command.
func NewSwitchCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "switch",
		Short:   "Switch account",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return switch_(c, cmd)
		},
	}

	cmd.Flags().StringVar(&switchOptions.account, "account", "", "Account name")
	return cmd
}

func switch_(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("account").Changed {
		switchOptions.account = c.Console().GetInput("username or organization")
	}
	conn, err := c.ClientConn()
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.SwitchRequest{
		Account: switchOptions.account,
	}
	header := metadata.MD{}
	_, err = client.Switch(context.Background(), request, grpc.Header(&header))
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if err := cli.SaveToken(header); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Printf("You are now logged in as: %s\n", switchOptions.account)
	return nil
}
