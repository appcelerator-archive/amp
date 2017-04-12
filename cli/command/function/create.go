package function_

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/function"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type createFunctionOpts struct {
	name  string
	image string
}

// NewFunctionCreateCommand returns a new instance of the function create command.
func NewFunctionCreateCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "create FUNCTION IMAGE",
		Short:   "Create function",
		PreRunE: cli.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := &createFunctionOpts{}
			opts.name = args[0]
			opts.image = args[1]
			return createFunction(c, opts)
		},
	}
}

func createFunction(c cli.Interface, opts *createFunctionOpts) error {
	client := function.NewFunctionClient(c.ClientConn())
	request := &function.CreateRequest{
		Name:  opts.name,
		Image: opts.image,
	}
	reply, err := client.Create(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println(reply.Function.Id)
	return nil
}
