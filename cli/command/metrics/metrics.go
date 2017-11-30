package metrics

import (
	"context"
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/metrics"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

type metricsOptions struct {
	cpu       bool
	mem       bool
	disk      bool
	net       bool
	timeRange string
	average   bool
}

// NewMetricsCommand returns a new instance of the metrics command.
func NewMetricsCommand(c cli.Interface) *cobra.Command {
	opts := metricsOptions{}
	cmd := &cobra.Command{
		Use:     "metrics",
		Short:   "Metrics management operations",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getMetrics(c, opts)
		},
	}

	//cmd.AddCommand(NewCPUCommand(c))
	flags := cmd.Flags()
	// Metrics options
	flags.BoolVar(&opts.cpu, "cpu", false, "Display CPU usage")
	flags.BoolVar(&opts.mem, "mem", false, "Display Memory usage")
	flags.BoolVar(&opts.disk, "disk", false, "Display Disk I/O usage")
	flags.BoolVar(&opts.net, "net", false, "Display Network I/O usage")
	// Query Options
	flags.StringVarP(&opts.timeRange, "time", "t", "1", "Time Range Duration")
	flags.BoolVar(&opts.average, "average", false, "Display average CPU usage")

	return cmd
}

func getMetrics(c cli.Interface, opts metricsOptions) error {
	conn := c.ClientConn()
	client := metrics.NewMetricsClient(conn)
	request := &metrics.MetricsRequest{
		Cpu:       opts.cpu,
		Mem:       opts.mem,
		Disk:      opts.disk,
		Net:       opts.disk,
		TimeRange: opts.timeRange,
		Average:   opts.average,
	}
	response, err := client.MetricsQuery(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
		return fmt.Errorf("error querying metrics: %s", err)
	}
	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "SERVICE NAME\tCPU USAGE")
	for _, entry := range response.Entries {
		fmt.Fprintf(w, "%s\t%.3f", entry.Service, entry.Usage)
	}
	return nil
}
