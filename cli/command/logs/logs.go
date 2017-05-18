package logs

import (
	"context"
	"fmt"
	"io"

	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type logsOptions struct {
	follow         bool
	includeAmpLogs bool
	meta           bool
	number         int64
	msg            string
	container      string
	stack          string
	node           string
}

// NewLogsCommand returns a new instance of the logs command.
func NewLogsCommand(c cli.Interface) *cobra.Command {
	opts := logsOptions{}
	cmd := &cobra.Command{
		Use:   "logs [OPTIONS] SERVICE",
		Short: "Fetch log entries matching provided criteria",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getLogs(c, args, opts)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&opts.follow, "follow", "f", false, "Follow log output")
	flags.BoolVarP(&opts.includeAmpLogs, "include", "i", false, "Include AMP logs")
	flags.BoolVarP(&opts.meta, "meta", "m", false, "Display entry metadata")
	flags.Int64VarP(&opts.number, "number", "n", 1000, "Number of results")
	flags.StringVar(&opts.msg, "msg", "", "Filter the message content by the given pattern")
	flags.StringVar(&opts.container, "container", "", "Filter by the given container")
	flags.StringVar(&opts.stack, "stack", "", "Filter by the given stack")
	flags.StringVar(&opts.node, "node", "", "Filter by the given node")
	return cmd
}

func getLogs(c cli.Interface, args []string, opts logsOptions) error {
	request := logs.GetRequest{}
	if len(args) > 0 {
		request.Service = args[0]
	}
	request.Message = opts.msg
	request.Container = opts.container
	request.Stack = opts.stack
	request.Node = opts.node
	request.Size = opts.number
	request.IncludeAmpLogs = opts.includeAmpLogs

	// Get logs from amplifier
	ctx := context.Background()
	conn := c.ClientConn()
	lc := logs.NewLogsClient(conn)
	r, err := lc.Get(ctx, &request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	for _, entry := range r.Entries {
		displayLogEntry(c, entry, opts.meta)
	}
	if !opts.follow {
		return nil
	}

	// If follow is requested, get subsequent logs and stream it
	stream, err := lc.GetStream(ctx, &request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	for {
		entry, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("%s", grpc.ErrorDesc(err))
		}
		displayLogEntry(c, entry, opts.meta)
	}
	return nil
}

func displayLogEntry(c cli.Interface, entry *logs.LogEntry, meta bool) {
	if meta {
		c.Console().Printf("%+v\n", entry)
	} else {
		c.Console().Printf("%s\n", entry.Msg)
	}
}
