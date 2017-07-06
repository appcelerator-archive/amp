package service

import (
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/cli/command/logs"
	"github.com/spf13/cobra"
)

// NewServiceLogsCommand returns a new instance of the service logs command.
func NewServiceLogsCommand(c cli.Interface) *cobra.Command {
	opts := logs.LogsOptions{}
	cmd := &cobra.Command{
		Use:     "logs [OPTIONS] SERVICE",
		Short:   "Fetch log entries of given service matching provided criteria",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logs.GetLogs(c, args, opts)
		},
	}
	flags := cmd.Flags()
	logs.AddLogFlags(flags, &opts)
	return cmd
}
