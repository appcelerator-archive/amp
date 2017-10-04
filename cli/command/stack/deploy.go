package stack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/docker/cli/cli/config"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

type deployStackOptions struct {
	file         string
	envVar       string
	registryAuth bool
}

var (
	opts = deployStackOptions{}
)

// NewDeployCommand returns a new instance of the stack command.
func NewDeployCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deploy [OPTIONS] STACK",
		Aliases: []string{"up", "start"},
		Short:   "Deploy a stack with a Docker Compose v3 file",
		PreRunE: cli.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy(c, cmd, args)
		},
	}
	cmd.Flags().StringVarP(&opts.file, "compose-file", "c", "", "Path to a Compose v3 file")
	cmd.Flags().StringVarP(&opts.envVar, "env", "e", "", "Environment variable to set during deployment format: var=value")
	cmd.Flags().BoolVar(&opts.registryAuth, "with-registry-auth", false, "Send registry authentication details to swarm agents")
	return cmd
}

func deploy(c cli.Interface, cmd *cobra.Command, args []string) error {
	envArgs := make(map[string]string)
	if opts.envVar != "" {
		envs := strings.Split(opts.envVar, "=")
		if len(envs) != 2 {
			return fmt.Errorf("--env parameter format error, should be: var=value found:: %s", opts.envVar)
		}
		envArgs[envs[0]] = envs[1]
	}
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
		EnvVar:  envArgs,
	}

	// If registryAuth was set, send the docker CLI configuration file to amplifier
	if opts.registryAuth {
		cf, err := config.Load(config.Dir())
		if err != nil {
			return errors.New("Unable to read docker CLI configuration file.")
		}
		req.Config, err = json.Marshal(cf)
		if err != nil {
			return errors.New("Unable to marshal docker CLI configuration file.")
		}
	}

	client := stack.NewStackClient(c.ClientConn())
	reply, err := client.Deploy(context.Background(), req)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	c.Console().Print(reply.Answer)
	return nil
}
