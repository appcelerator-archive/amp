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

type listOpts struct {
	quiet bool
}

var (
	lsopts = &listOpts{}
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
	lines := strings.Split(reply.Answer, "\n")
	if len(lines) > 1 {
		if !lsopts.quiet {
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
		} else {
			for _, line := range lines[1:] {
				if line != "" {
					c.Console().Println(strings.Split(line, " ")[0])
				}
			}
		}

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
