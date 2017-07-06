package logs

import (
	"context"
	"errors"
	"io"
	"strconv"

	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/grpc/status"
)

type LogsOptions struct {
	Follow         bool
	IncludeAmpLogs bool
	Meta           bool
	Raw            bool
	Number         int64
	Msg            string
	Container      string
	Stack          string
	Node           string
}

func AddLogFlags(flags *pflag.FlagSet, opts *LogsOptions) {
	flags.BoolVarP(&opts.Follow, "follow", "f", false, "Follow log output")
	flags.BoolVarP(&opts.IncludeAmpLogs, "include", "i", false, "Include AMP logs")
	flags.BoolVarP(&opts.Meta, "meta", "m", false, "Display entry metadata")
	flags.Int64VarP(&opts.Number, "number", "n", 1000, "Number of results")
	flags.StringVar(&opts.Msg, "msg", "", "Filter the message content by the given pattern")
	flags.StringVar(&opts.Container, "container", "", "Filter by the given Container")
	flags.StringVar(&opts.Node, "node", "", "Filter by the given node")
	flags.BoolVarP(&opts.Raw, "raw", "r", false, "Display raw logs (no prefix)")
}

// NewLogsCommand returns a new instance of the logs command.
func NewLogsCommand(c cli.Interface) *cobra.Command {
	opts := LogsOptions{}
	cmd := &cobra.Command{
		Use:   "logs [OPTIONS] SERVICE",
		Short: "Fetch log entries matching provided criteria",
		RunE: func(cmd *cobra.Command, args []string) error {
			return GetLogs(c, args, opts)
		},
	}
	flags := cmd.Flags()
	AddLogFlags(flags, &opts)
	flags.StringVar(&opts.Stack, "stack", "", "Filter by the given stack")
	return cmd
}

func GetLogs(c cli.Interface, args []string, opts LogsOptions) error {
	request := logs.GetRequest{}
	if len(args) > 0 {
		request.Service = args[0]
	}
	request.Message = opts.Msg
	request.Container = opts.Container
	request.Stack = opts.Stack
	request.Node = opts.Node
	request.Size = opts.Number
	request.IncludeAmpLogs = opts.IncludeAmpLogs

	// Get logs from amplifier
	ctx := context.Background()
	conn := c.ClientConn()
	lc := logs.NewLogsClient(conn)
	r, err := lc.Get(ctx, &request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	for _, entry := range r.Entries {
		displayLogEntry(c, entry, opts.Meta, opts.Raw)
	}
	if !opts.Follow {
		return nil
	}

	// If Follow is requested, get subsequent logs and stream it
	stream, err := lc.GetStream(ctx, &request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	for {
		entry, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			if s, ok := status.FromError(err); ok {
				return errors.New(s.Message())
			}
		}
		displayLogEntry(c, entry, opts.Meta, opts.Raw)
	}
	return nil
}

func displayLogEntry(c cli.Interface, entry *logs.LogEntry, meta bool, raw bool) {
	if meta {
		c.Console().Printf("%+v\n", entry)
	} else if raw {
		c.Console().Printf("%s\n", entry.Msg)
	} else {
		c.Console().Printf("%24s | %s\n", entry.ServiceName+"."+strconv.Itoa(int(entry.TaskSlot)), entry.Msg)
	}
}
