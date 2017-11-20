package stack

import (
	"context"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

type listStackOptions struct {
	quiet bool
}

// NewListCommand returns a new instance of the stack command.
func NewListCommand(c cli.Interface) *cobra.Command {
	opts := listStackOptions{}
	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Short:   "List deployed stacks",
		Aliases: []string{"list"},
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(c, opts)
		},
	}
	cmd.Flags().BoolVarP(&opts.quiet, "quiet", "q", false, "Only display stack ids")
	return cmd
}

func list(c cli.Interface, opts listStackOptions) error {
	req := &stack.ListRequest{}
	client := stack.NewStackClient(c.ClientConn())
	reply, err := client.StackList(context.Background(), req)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	if opts.quiet {
		for _, line := range reply.Entries {
			c.Console().Println(line.Stack.Id)
		}
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tRUNNING\tCOMPLETE\tPREPARING\tTOTAL\tSERVICES\tSTATUS\tOWNER")
	for _, entry := range reply.Entries {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%d\t%d\t%d/%d\t%s\t%s\n", entry.Stack.Id, entry.Stack.Name, entry.RunningServices, entry.CompleteServices, entry.PreparingServices, entry.TotalServices, entry.RunningServices, entry.TotalServices, entry.Status, entry.Stack.Owner.User)
	}
	w.Flush()
	return nil
}
