package stack

import (
	"errors"
	"io/ioutil"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type deployOpts struct {
	name string
	file string
}

var (
	opts = &deployOpts{}
)

// NewDeployCommand returns a new instance of the stack command.
func NewDeployCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deploy",
		Short:   "Deploy a stack with a docker compose v3 file",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.name = args[0]
			return deploy(c)
		},
	}
	cmd.Flags().StringVarP(&opts.file, "compose-file", "c", "", "Path to a Compose v3 file")
	return cmd
}

func deploy(c cli.Interface) error {
	c.Console().Printf("Deploying stack %s using %s\n", opts.name, opts.file)

	contents, err := ioutil.ReadFile(opts.file)
	if err != nil {
		return err
	}

	req := &stack.DeployRequest{
		Name:    opts.name,
		Compose: contents,
	}

	client := stack.NewStackClient(c.ClientConn())
	reply, err := client.Deploy(context.Background(), req)
	if err != nil {
		return errors.New(grpc.ErrorDesc(err))
	}
	c.Console().Println(reply.Answer)
	return nil
}
