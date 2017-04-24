package org

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// NewSwitchCommand returns a new instance of the switch command.
func NewSwitchCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "switch ACCOUNT",
		Short:   "Switch account",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return switch_(c, args)
		},
	}
}

func switch_(c cli.Interface, args []string) error {
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.SwitchRequest{
		Account: args[0],
	}
	header := metadata.MD{}
	_, err := client.Switch(context.Background(), request, grpc.Header(&header))
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if err := cli.SaveToken(header); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Printf("You are now logged in as: %s\n", args[0])
	return nil
}
