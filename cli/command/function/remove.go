package function_

import (
	"fmt"

	"errors"

	"github.com/appcelerator/amp/api/rpc/function"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type removeFunctionOpts struct {
	function string
}

// NewFunctionRemoveCommand returns a new instance of the function remove command.
func NewFunctionRemoveCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "rm FUNCTION",
		Short:   "Remove function",
		Aliases: []string{"del"},
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return errors.New("Function cannot be empty")
			}
			opts := &removeFunctionOpts{}
			opts.function = args[0]
			return removeFunction(c, opts)
		},
	}
}

func removeFunction(c cli.Interface, opts *removeFunctionOpts) error {
	client := function.NewFunctionClient(c.ClientConn())
	request := &function.DeleteRequest{
		Id: opts.function,
	}
	if _, err := client.Delete(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println(opts.function)
	return nil
}
