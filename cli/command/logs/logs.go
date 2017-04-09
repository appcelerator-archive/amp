package logs

import (
	"context"
	"io"

	"fmt"

	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type logsOpts struct {
	follow    bool
	infra     bool
	meta      bool
	number    int64
	msg       string
	container string
	stack     string
	node      string
}

var (
	opts = &logsOpts{}
)

// NewLogsCommand returns a new instance of the logs command.
func NewLogsCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs SERVICE",
		Short: "Fetch log entries matching provided criteria",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getLogs(c, args)
		},
	}
	cmd.Flags().BoolVarP(&opts.follow, "follow", "f", false, "Follow log output")
	cmd.Flags().BoolVarP(&opts.infra, "infra", "i", false, "Include infrastructure logs")
	cmd.Flags().BoolVarP(&opts.meta, "meta", "m", false, "Display entry metadata")
	cmd.Flags().Int64VarP(&opts.number, "number", "n", 1000, "Number of results")
	cmd.Flags().StringVar(&opts.msg, "msg", "", "Filter the message content by the given pattern")
	cmd.Flags().StringVar(&opts.container, "container", "", "Filter by the given container")
	cmd.Flags().StringVar(&opts.stack, "stack", "", "Filter by the given stack")
	cmd.Flags().StringVar(&opts.node, "node", "", "Filter by the given node")
	return cmd
}

func getLogs(c cli.Interface, args []string) error {
	request := logs.GetRequest{}
	if len(args) > 0 {
		request.Service = args[0]
	}
	request.Message = opts.msg
	request.Container = opts.container
	request.Stack = opts.stack
	request.Node = opts.node
	request.Size = opts.number
	request.Infra = opts.infra

	// Get logs from amplifier
	ctx := context.Background()
	conn, err := c.ClientConn()
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
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
