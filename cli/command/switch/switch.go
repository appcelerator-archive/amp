package switch_

import (
	"errors"
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
	return &cobra.Command{
		Use:     "switch ACCOUNT",
		Short:   "Switch account",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return errors.New("account name cannot be empty")
			}
			switchOptions.account = args[0]
			return switch_(c, switchOptions)
		},
	}
}

func switch_(c cli.Interface, opt *switchOpts) error {
	conn, err := c.ClientConn()
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.SwitchRequest{
		Account: opt.account,
	}
	header := metadata.MD{}
	_, err = client.Switch(context.Background(), request, grpc.Header(&header))
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if err := cli.SaveToken(header); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Printf("You are now logged in as: %s\n", opt.account)
	return nil
}
