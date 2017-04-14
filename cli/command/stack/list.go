package stack

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// NewListCommand returns a new instance of the stack command.
func NewListCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List deployed stacks",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(c)
		},
	}
	return cmd
}

func list(c cli.Interface) error {
	req := &stack.ListRequest{}
	client := stack.NewStackClient(c.ClientConn())
	reply, err := client.List(context.Background(), req)
	if err != nil {
		return errors.New(grpc.ErrorDesc(err))
	}
	lines := strings.Split(reply.Answer, "\n")
	if len(lines) > 1 {
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSERVICE")
		lines = lines[1:]
		sort.Strings(lines)
		for _, line := range lines {
			if line != "" {
				fmt.Fprintln(w, getOneStackListLine(line))
			}
		}
		w.Flush()
	}
	return nil
}

func getOneStackListLine(line string) string {
	cols := strings.Split(line, " ")
	name := cols[0]
	ll := strings.LastIndex(cols[0], "-")
	if ll >= 0 {
		name = name[0:ll]
	}
	ret := fmt.Sprintf("%s\t%s", cols[0], name)
	for _, col := range cols[1:] {
		if col != "" {
			ret = fmt.Sprintf("%s\t%s", ret, col)
		}
	}
	return ret
}
