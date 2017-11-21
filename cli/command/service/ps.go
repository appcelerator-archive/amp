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
		Use:     "ps [OPTIONS] SERVICE",
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
	request := &service.PsRequest{
		Service: args[0],
	}
	reply, err := client.ServicePs(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}

	prevName := ""

	w := tabwriter.NewWriter(os.Stdout, 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tIMAGE\tDESIRED STATE\tCURRENT STATE\tNODE ID\tERROR")

	for _, task := range reply.Tasks {
		if opts.quiet {
			c.Console().Println(task.Id)
			return nil
		} else {
			var name string
			if task.Slot != 0 {
				name = fmt.Sprintf("%v.%v", args[0], task.Slot)
			} else {
				name = fmt.Sprintf("%v.%v", args[0], task.NodeId)
			}
			// Indent the name if necessary
			indentedName := name
			if name == prevName {
				indentedName = fmt.Sprintf(" \\_ %s", indentedName)
			}
			prevName = name
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", task.Id, indentedName, task.Image, task.DesiredState, task.CurrentState, task.NodeId, task.Error)
		}
	}
	if !opts.quiet {
		w.Flush()
	}
	return nil
}
