package stack

import (
	"context"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type servicesOpts struct {
	quiet bool
}

var (
	sopts = &servicesOpts{}
)

// NewServicesCommand returns a new instance of the stack command.
func NewServicesCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "services STACK [OPTIONS]",
		Aliases: []string{"srv"},
		Short:   "List services of a stack",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return services(c, args[0])
		},
	}
	cmd.Flags().BoolVarP(&sopts.quiet, "quiet", "q", false, "Only display the stack id")
	return cmd
}

func services(c cli.Interface, stackName string) error {
	req := &stack.ServicesRequest{StackName: stackName}
	client := stack.NewStackClient(c.ClientConn())
	reply, err := client.Services(context.Background(), req)
	if err != nil {
		return errors.New(grpc.ErrorDesc(err))
	}
	if !sopts.quiet {
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tMODE\tREPLICAS\tIMAGE")
		for _, line := range reply.Services {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", line.Id, line.Name, line.Mode, line.Replicas, line.Image)
		}
		w.Flush()
	} else {
		for _, line := range reply.Services {
			c.Console().Println(line.Id)
		}
	}
	return nil
}
