package stack

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/docker/cli/cli/config"
	"github.com/appcelerator/amp/docker/cli/opts"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

type deployStackOptions struct {
	file         string
	env          opts.ListOpts
	envFile      opts.ListOpts
	registryAuth bool
}

var (
	options = deployStackOptions{
		env:     opts.NewListOpts(opts.ValidateEnv),
		envFile: opts.NewListOpts(nil),
	}
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
	cmd.Flags().StringVarP(&options.file, "compose-file", "c", "", "Path to a Compose v3 file")
	cmd.Flags().VarP(&options.env, "env", "e", "Set environment variables")
	cmd.Flags().Var(&options.envFile, "env-file", "Read in a file of environment variables")
	cmd.Flags().BoolVar(&options.registryAuth, "with-registry-auth", false, "Send registry authentication details to swarm agents")
	return cmd
}

func deploy(c cli.Interface, cmd *cobra.Command, args []string) error {
	// Environment
	environment, err := opts.ReadKVStrings(options.envFile.GetAll(), options.env.GetAll())
	if err != nil {
		return err
	}
	currentEnv := make([]string, 0, len(environment))
	for _, env := range environment { // need to process each var, in order
		k := strings.SplitN(env, "=", 2)[0]
		for i, current := range currentEnv { // remove duplicates
			if current == env {
				continue // no update required, may hide this behind flag to preserve order of environment
			}
			if strings.HasPrefix(current, k+"=") {
				currentEnv = append(currentEnv[:i], currentEnv[i+1:]...)
			}
		}
		currentEnv = append(currentEnv, env)
	}

	// Compose file
	var name string
	if len(args) == 0 {
		basename := filepath.Base(options.file)
		name = strings.Split(strings.TrimSuffix(basename, filepath.Ext(options.file)), ".")[0]
	} else {
		name = args[0]
	}
	c.Console().Printf("Deploying stack %s using %s\n", name, options.file)

	contents, err := ioutil.ReadFile(options.file)
	if err != nil {
		return err
	}

	req := &stack.DeployRequest{
		Name:        name,
		Compose:     contents,
		Environment: environment,
	}

	// If registryAuth was set, send the docker CLI configuration file to amplifier
	if options.registryAuth {
		cf, err := config.Load(config.Dir())
		if err != nil {
			return errors.New("unable to read docker CLI configuration file")
		}
		req.Config, err = json.Marshal(cf)
		if err != nil {
			return errors.New("unable to marshal docker CLI configuration file")
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
