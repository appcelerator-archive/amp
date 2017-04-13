package stack

import (
	"context"
	"errors"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type removeOpts struct {
	name string
}

var (
	ropts = &removeOpts{}
)

// NewRemoveCommand returns a new instance of the stack command.
func NewRemoveCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove STACKNAME",
		Aliases: []string{"rm", "down", "stop"},
		Short:   "Remove a deployed stack",
		RunE: func(cmd *cobra.Command, args []string) error {
			ropts.name = args[0]
			return remove(c)
		},
	}
	return cmd
}

func remove(c cli.Interface) error {
	req := &stack.RemoveRequest{
		Id: ropts.name,
	}

	client := stack.NewStackClient(c.ClientConn())
	reply, err := client.Remove(context.Background(), req)
	if err != nil {
		return errors.New(grpc.ErrorDesc(err))
	}
	c.Console().Println(reply.Answer)
	return nil
}
