package config

import (
	"errors"
	"fmt"

	"github.com/appcelerator/amp/api/rpc/config"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

// NewRemoveCommand returns a new instance of the remove command for removing one or more configs
func NewRemoveCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove [OPTIONS]",
		Short:   "Remove one or more configs",
		Aliases: []string{"rm"},
		PreRunE: cli.AtLeastArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return remove(c, cmd, args)
		},
	}

	return cmd
}

func remove(c cli.Interface, cmd *cobra.Command, args []string) error {
	conn := c.ClientConn()
	client := config.NewConfigClient(conn)
	for _, id := range args {
		request := &config.RemoveRequest{Id: id}
		reply, err := client.ConfigRemove(context.Background(), request)
		if err != nil {
			if s, ok := status.FromError(err); ok {
				return errors.New(s.Message())
			}
			return fmt.Errorf("error removing config: %s", err)
		}
		c.Console().Println(reply.GetId())
	}
	return nil
}
