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
	"google.golang.org/grpc"
)

type listOpts struct {
	quiet bool
}

var (
	lsopts = &listOpts{}
)

// NewListCommand returns a new instance of the stack command.
func NewListCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "List deployed stacks",
		Aliases: []string{"list"},
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(c)
		},
	}
	cmd.Flags().BoolVarP(&lsopts.quiet, "quiet", "q", false, "Only display the stack id")
	return cmd
}

func list(c cli.Interface) error {
	req := &stack.ListRequest{}
	client := stack.NewStackClient(c.ClientConn())
	reply, err := client.List(context.Background(), req)
	if err != nil {
		return errors.New(grpc.ErrorDesc(err))
	}
	if !lsopts.quiet {
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSERVICE\tOWNER\tOWNER TYPE")
		for _, entry := range reply.Entries {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", entry.Stack.Id, entry.Stack.Name, entry.Services, entry.Stack.Owner.Name, entry.Stack.Owner.Type)
		}
		w.Flush()
	} else {
		for _, line := range reply.Entries {
			c.Console().Println(line.Stack.Id)
		}
	}
	return nil
}
