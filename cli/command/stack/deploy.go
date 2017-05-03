package stack

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type deployStackOptions struct {
	file string
}

var (
	opts = deployStackOptions{}
)

// NewDeployCommand returns a new instance of the stack command.
func NewDeployCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deploy [OPTIONS] STACK",
		Aliases: []string{"up", "start"},
		Short:   "Deploy a stack with a docker compose v3 file",
		PreRunE: cli.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy(c, args)
		},
	}
	cmd.Flags().StringVarP(&opts.file, "compose-file", "c", "", "Path to a Compose v3 file")
	return cmd
}

func deploy(c cli.Interface, args []string) error {
	var name string
	if len(args) == 0 {
		basename := filepath.Base(opts.file)
		name = strings.Split(strings.TrimSuffix(basename, filepath.Ext(opts.file)), ".")[0]
	} else {
		name = args[0]
	}
	c.Console().Printf("Deploying stack %s using %s\n", name, opts.file)

	contents, err := ioutil.ReadFile(opts.file)
	if err != nil {
		return err
	}

	req := &stack.DeployRequest{
		Name:    name,
		Compose: contents,
	}

	client := stack.NewStackClient(c.ClientConn())
	reply, err := client.Deploy(context.Background(), req)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println(reply.Answer)
	return nil
}
