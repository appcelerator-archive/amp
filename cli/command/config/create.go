package config

import (
	"context"
	"errors"

	"github.com/appcelerator/amp/api/rpc/config"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

type CreateOpts struct {
	Labels []string
}

var createOpts = &CreateOpts{
	Labels: []string{},
}

// NewCreateCommand returns a new instance of the create command for creating a config.
func NewCreateCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [OPTIONS]",
		Short:   "Create a config",
		PreRunE: cli.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(c, cmd)
		},
	}

	flags := cmd.Flags()
	flags.StringSliceVarP(&createOpts.Labels, "labels", "l", []string{}, "Secret labels")

	return cmd
}

func create(c cli.Interface, cmd *cobra.Command) error {
	conn := c.ClientConn()
	client := config.NewConfigClient(conn)
	request := &config.CreateConfigRequest{}
	if _, err := client.CreateConfig(context.Background(), request); err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	c.Console().Println("Configuration successfully created.")
	return nil
}
