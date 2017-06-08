package service

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// NewServiceInspectCommand returns a new instance of the service inspect command.
func NewServiceInspectCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "inspect SERVICE",
		Short:   "Display detailed information of a service",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return inspectService(c, args)
		},
	}
}

func inspectService(c cli.Interface, args []string) error {
	conn := c.ClientConn()
	client := service.NewServiceClient(conn)
	request := &service.ServiceInspectRequest{
		ServiceId: args[0],
	}
	reply, err := client.InspectService(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println(reply.ServiceEntity)
	return nil
}
