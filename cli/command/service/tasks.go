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
	"google.golang.org/grpc"
)

type tasksOpts struct {
	quiet bool
}

var (
	topts = &tasksOpts{}
)

// NewServiceTasksCommand returns a new instance of the stack command.
func NewServiceTasksCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tasks",
		Short:   "List the service tasks",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return tasks(c, args[0])
		},
	}
	cmd.Flags().BoolVarP(&topts.quiet, "quiet", "q", false, "Only display the task id")
	return cmd
}

func tasks(c cli.Interface, serviceID string) error {
	req := &service.TasksRequest{ServiceId: serviceID}
	client := service.NewServiceClient(c.ClientConn())
	reply, err := client.Tasks(context.Background(), req)
	if err != nil {
		return errors.New(grpc.ErrorDesc(err))
	}
	if !topts.quiet {
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tIMAGE\tDESIRED STATE\tSTATE\tNODEID")
		for _, line := range reply.Tasks {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", line.Id, line.Image, line.DesiredState, line.State, line.NodeId)
		}
		w.Flush()
	} else {
		for _, line := range reply.Tasks {
			c.Console().Println(line.Id)
		}
	}
	return nil
}
