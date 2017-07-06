package stack

import (
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/cli/command/logs"
	"github.com/spf13/cobra"
)

type logsStackOptions struct {
	meta   bool
	follow bool
}

// NewLogsCommand returns a new instance of the stack command.
func NewLogsCommand(c cli.Interface) *cobra.Command {
	opts := logs.LogsOptions{}
	cmd := &cobra.Command{
		Use:     "logs [OPTIONS] STACK",
		Short:   "Fetch log entries of given stack matching provided criteria",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Stack = args[0]
			args = nil
			return logs.GetLogs(c, args, opts)
		},
	}
	flags := cmd.Flags()
	logs.AddLogFlags(flags, &opts)
	return cmd
}
