package metrics

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

type metricsOptions struct {
	cpu bool
	mem bool
	disk bool
	net bool
	timeRange string
	average bool
}

// NewMetricsCommand returns a new instance of the metrics command.
func NewMetricsCommand(c cli.Interface) *cobra.Command {
	opts := metricsOptions{}
	cmd := &cobra.Command{
		Use:     "metrics",
		Short:   "Metrics management operations",
		PreRunE: cli.NoArgs,
		RunE:    func(cmd *cobra.Command, args []string) error {
			return metrics(c, opts)
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
	flags.BoolVar(&opts.average, "average",  false, "Display average CPU usage")

	return cmd
}

func metrics(c cli.Interface, opts metricsOptions) error {
	conn := c.ClientConn()
	client := metrics(conn)
	request := &metrics.MetricsRequest{

	}

	return nil
}