package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"docker.io/go-docker/api/types/swarm"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/docker/cli/cli/command/formatter"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

// inspectOptions as defined in https://github.com/docker/cli/blob/master/cli/command/service/inspect.go
type inspectOptions struct {
	refs   []string
	format string
	pretty bool
}

var opts inspectOptions

// NewServiceInspectCommand returns a new instance of the service inspect command.
func NewServiceInspectCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inspect SERVICE",
		Short:   "Display detailed information of a service",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return inspectService(c, args)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&opts.format, "format", "f", "", "Format the output using the given Go template")
	return cmd
}

func inspectService(c cli.Interface, args []string) error {
	conn := c.ClientConn()
	client := service.NewServiceClient(conn)
	request := &service.InspectRequest{
		Service: args[0],
	}
	reply, err := client.ServiceInspect(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	if len(opts.format) == 0 {
		c.Console().Println(reply.Json)
		return nil
	}
	opts.refs = args
	// check if the user is trying to apply a template to the pretty format, which
	// is not supported
	if strings.HasPrefix(opts.format, "pretty") && opts.format != "pretty" {
		return errors.New("cannot supply extra formatting options to the pretty template")
	}
	serviceCtx := formatter.Context{
		Output: c.Console().OutStream(),
		Format: formatter.NewServiceFormat(opts.format),
	}
	getRef := func(ref string) (interface{}, []byte, error) {
		var s swarm.Service
		r := bytes.NewReader([]byte(reply.Json))
		err := json.NewDecoder(r).Decode(&s)
		if err != nil {
			return nil, nil, err
		}
		return s, nil, nil
	}
	if err := formatter.ServiceInspectWrite(serviceCtx, opts.refs, getRef, nil); err != nil {
		return err
	}
	return nil
}
