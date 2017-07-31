package config

import (
	"context"
	"errors"

	"strings"

	"github.com/appcelerator/amp/api/rpc/config"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

// NewRemoveCommand returns a new instance of the remove command for removing one or more configs
func NewRemoveCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "rm [OPTIONS]",
		Short:   "Remove one or more configs",
		Aliases: []string{"remove"},
		PreRunE: cli.AtLeastArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return remove(c, args)
		},
	}
}

func remove(c cli.Interface, args []string) error {
	var errs []string
	conn := c.ClientConn()
	client := config.NewConfigClient(conn)
	for _, cfg := range args {
		request := &config.RemoveConfigRequest{
			Name: cfg,
		}
		if _, err := client.RemoveConfig(context.Background(), request); err != nil {
			if s, ok := status.FromError(err); ok {
				errs = append(errs, s.Message())
				continue
			}
		}
		c.Console().Println(cfg)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
