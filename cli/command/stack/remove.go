package stack

import (
	"errors"
	"strings"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// NewRemoveCommand returns a new instance of the stack command.
func NewRemoveCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "rm STACK(S)",
		Aliases: []string{"remove", "down", "stop"},
		Short:   "Remove a deployed stack",
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
		if _, err := client.Remove(context.Background(), req); err != nil {
			errs = append(errs, grpc.ErrorDesc(err))
			continue
		}
		c.Console().Println(name)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
