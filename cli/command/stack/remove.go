package stack

import (
	"errors"
	"strings"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

// NewRemoveCommand returns a new instance of the stack command.
func NewRemoveCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "rm STACK(S)",
		Aliases: []string{"remove", "down", "stop"},
		Short:   "Remove one or more deployed stacks",
		PreRunE: cli.AtLeastArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return remove(c, args)
		},
	}
}

func remove(c cli.Interface, args []string) error {
	var errs []string
	conn := c.ClientConn()
	client := stack.NewStackClient(conn)
	for _, name := range args {
		req := &stack.RemoveRequest{
			Stack: name,
		}
		reply, err := client.StackRemove(context.Background(), req)
		if err != nil {
			if s, ok := status.FromError(err); ok {
				errs = append(errs, s.Message())
				continue
			}
		}
		c.Console().Print(reply.Answer)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
