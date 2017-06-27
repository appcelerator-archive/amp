package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

type taskOptions struct {
	quiet bool
}

// NewServicePsCommand returns a new instance of the service ps command
func NewServicePsCommand(c cli.Interface) *cobra.Command {
	opts := taskOptions{}
	cmd := &cobra.Command{
		Use:     "ps SERVICE [OPTIONS]",
		Short:   "List tasks of a service",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return tasks(c, args, opts)
		},
	}
	cmd.Flags().BoolVarP(&opts.quiet, "quiet", "q", false, "Only display task ids")
	return cmd
}

func tasks(c cli.Interface, args []string, opts taskOptions) error {
	conn := c.ClientConn()
	client := service.NewServiceClient(conn)
	request := &service.TasksRequest{
		ServiceId: args[0],
	}
	reply, err := client.Tasks(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	if opts.quiet {
		for _, task := range reply.Tasks {
			c.Console().Println(task.Id)
		}
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "ID\tIMAGE\tDESIRED STATE\tCURRENT STATE\tNODE ID\tERROR")
	for _, task := range reply.Tasks {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", task.Id, task.Image, task.DesiredState, task.CurrentState, task.NodeId, task.Error)
	}
	w.Flush()
	return nil
}
