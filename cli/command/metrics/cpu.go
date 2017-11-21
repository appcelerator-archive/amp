package metrics

import (
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/metrics"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

type cpuOptions struct {
	average bool
}

// NewCPUCommand returns a new instance of the cpu command.
func NewCPUCommand(c cli.Interface) *cobra.Command {
	opts := cpuOptions{}
	cmd := &cobra.Command{
		Use:     "cpu [OPTIONS]",
		Short:   "Display metrics for Container CPU usage",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cpu(c, opts)
		},
	}
	cmd.Flags().BoolVar(&opts.average, "average", false, "Average CPU Usage")
	return cmd
}

func cpu(c cli.Interface, opts cpuOptions) error {
	conn := c.ClientConn()
	client := metrics.NewMetricsClient(conn)
	request := &metrics.CPUMetricsRequest{
		Average: opts.average,
	}
	response, err := client.CPUMetricsQuery(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}

	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	// Tabulate the results
	fmt.Fprintln(w, "SERVICE NAME\tCPU USAGE (in %)")
	for _, entry := range response.Entries {
		fmt.Fprintf(w, "%s\t%.3f\n", entry.Service, entry.Usage)
	}
	w.Flush()
	return nil
}
