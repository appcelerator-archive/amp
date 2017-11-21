package metrics

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewMetricsCommand returns a new instance of the metrics command.
func NewMetricsCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "metrics",
		Short:   "Metrics management operations",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}

	cmd.AddCommand(NewCPUCommand(c))

	return cmd
}
