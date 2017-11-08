package service

import (
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

type listServiceOptions struct {
	quiet bool
	stack string
}

// NewServiceListCommand returns a new instance of the service list command.
func NewServiceListCommand(c cli.Interface) *cobra.Command {
	opts := listServiceOptions{}
	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Short:   "List services",
		Aliases: []string{"list"},
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listServices(c, opts)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only display service ids")
	flags.StringVar(&opts.stack, "stack", "", "Filter services based on stack name")
	return cmd
}

func listServices(c cli.Interface, opts listServiceOptions) error {
	conn := c.ClientConn()
	client := service.NewServiceClient(conn)
	request := &service.ListRequest{
		Stack: opts.stack,
	}
	reply, err := client.List(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	if opts.quiet {
		for _, entry := range reply.Entries {
			c.Console().Println(entry.Id)
		}
		return nil
	}
	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tMODE\tREPLICAS\tSTATUS\tIMAGE\tTAG")
	for _, entry := range reply.Entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d/%d\t%s\t%s\t%s\n", entry.Id, entry.Name, entry.Mode, entry.RunningTasks, entry.TotalTasks, entry.Status, entry.Image, entry.Tag)
	}
	w.Flush()
	return nil
}
