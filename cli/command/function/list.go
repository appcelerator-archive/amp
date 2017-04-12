package function_

import (
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/function"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/pkg/time"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type listFunctionOpts struct {
	quiet bool
}

var (
	listFunctionOptions = &listFunctionOpts{}
)

// NewFunctionListCommand returns a new instance of the function list command.
func NewFunctionListCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Short:   "List function",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listFunction(c)
		},
	}
	cmd.Flags().BoolVarP(&listFunctionOptions.quiet, "quiet", "q", false, "Only display function id")
	return cmd
}

func listFunction(c cli.Interface) error {
	client := function.NewFunctionClient(c.ClientConn())
	request := &function.ListRequest{}
	reply, err := client.List(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if listFunctionOptions.quiet {
		for _, f := range reply.Functions {
			c.Console().Println(f.Id)
		}
		return nil
	}
	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tIMAGE\tOWNER\tCREATED ON")
	for _, f := range reply.Functions {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", f.Id, f.Name, f.Image, f.Owner.Name, time.ConvertTime(f.CreateDt))
	}
	w.Flush()
	return nil
}
