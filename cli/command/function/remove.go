package function_

import (
	"errors"
	"strings"

	"github.com/appcelerator/amp/api/rpc/function"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// NewFunctionRemoveCommand returns a new instance of the function remove command.
func NewFunctionRemoveCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "rm FUNCTION(S)",
		Short:   "Remove function",
		Aliases: []string{"remove"},
		PreRunE: cli.AtLeastArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeFunction(c, args)
		},
	}
}

func removeFunction(c cli.Interface, args []string) error {
	var errs []string
	client := function.NewFunctionClient(c.ClientConn())
	for _, fn := range args {
		request := &function.DeleteRequest{
			Id: fn,
		}
		if _, err := client.Delete(context.Background(), request); err != nil {
			errs = append(errs, grpc.ErrorDesc(err))
			continue
		}
		c.Console().Println(fn)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
