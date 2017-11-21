package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

type scaleOptions struct {
	service  string
	replicas uint64
}

// NewServiceScaleCommand returns a new instance of the service scale command
func NewServiceScaleCommand(c cli.Interface) *cobra.Command {
	opts := scaleOptions{}
	cmd := &cobra.Command{
		Use:     "scale [OPTIONS]",
		Short:   "Scale a replicated service",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return scale(c, cmd, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.service, "service", "", "Service ID or name")
	flags.Uint64VarP(&opts.replicas, "replicas", "r", 0, "Number of replicas")

	return cmd
}

func scale(c cli.Interface, cmd *cobra.Command, opts scaleOptions) error {
	if !cmd.Flag("service").Changed {
		opts.service = c.Console().GetInput("service id or name")
	}
	if !cmd.Flag("replicas").Changed {
		rep, err := strconv.Atoi(c.Console().GetInput("replicas"))
		if err != nil {
			return fmt.Errorf("%s", err.Error())
		}
		opts.replicas = uint64(rep)
	}

	conn := c.ClientConn()
	client := service.NewServiceClient(conn)
	request := &service.ScaleRequest{
		Service:  opts.service,
		Replicas: opts.replicas,
	}
	if _, err := client.ServiceScale(context.Background(), request); err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	c.Console().Println("Service", opts.service, "has been scaled to", opts.replicas, "replicas.")
	return nil
}
