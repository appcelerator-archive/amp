package stack

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
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
	reply, err := client.List(context.Background(), req)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if opts.quiet {
		for _, line := range reply.Entries {
			c.Console().Println(line.Stack.Id)
		}
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tSERVICES\tFAILED SERVICES\tSTATUS\tOWNER\tORGANIZATION")
	for _, entry := range reply.Entries {
		fmt.Fprintf(w, "%s\t%s\t%d/%d\t%d\t%s\t%s\t%s\n", entry.Stack.Id, entry.Stack.Name, entry.RunningServices, entry.TotalServices, entry.FailedServices, entry.Status, entry.Stack.Owner.User, entry.Stack.Owner.Organization)
	}
	w.Flush()
	return nil
}
