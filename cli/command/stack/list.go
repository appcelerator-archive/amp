package stack

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// NewStackCommand returns a new instance of the stack command.
func NewListCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
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
	c.Console().Printf("%s\n", lines[0])
	if len(lines) > 1 {
		lines = lines[1:]
		sort.Strings(lines)
		for _, line := range lines {
			if line != "" {
				c.Console().Printf("%s\n", line)
			}
		}
	}
	return nil
}
