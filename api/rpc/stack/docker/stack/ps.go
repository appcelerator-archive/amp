package stack

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	//"github.com/docker/docker/cli" //{AMP}:remove it
	"github.com/docker/docker/cli/command"
	"github.com/docker/docker/cli/command/idresolver"
	"github.com/docker/docker/cli/command/task"
	"github.com/docker/docker/opts"
	//"github.com/spf13/cobra" {AMP}:remove it
)

type psOptions struct {
	filter    opts.FilterOpt
	noTrunc   bool
	namespace string
	noResolve bool
}

// NewPsOptions {AMP} Add public construct to this private strcut
func NewPsOptions(namespace string, noTrunc bool, noResolve bool, filter string) psOptions {
	optsf := opts.NewFilterOpt()
	optsf.Set(filter)
	return psOptions{
		namespace: namespace,
		noTrunc:   noTrunc,
		noResolve: noResolve,
		filter:    optsf,
	}
}

/* {AMP}: remove this function and use it on cli side
func NewPsCommand(dockerCli *command.DockerCli) *cobra.Command {
	opts := psOptions{filter: opts.NewFilterOpt()}

	cmd := &cobra.Command{
		Use:   "ps [OPTIONS] STACK",
		Short: "List the tasks in the stack",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.namespace = args[0]
			return runPS(dockerCli, opts)
		},
	}
	flags := cmd.Flags()
	flags.BoolVar(&opts.noTrunc, "no-trunc", false, "Do not truncate output")
	flags.BoolVar(&opts.noResolve, "no-resolve", false, "Do not map IDs to Names")
	flags.VarP(&opts.filter, "filter", "f", "Filter output based on conditions provided")

	return cmd
}
*/

// RunPS {AMP}:make it public
func RunPS(dockerCli *command.DockerCli, opts psOptions) error {
	namespace := opts.namespace
	client := dockerCli.Client()
	ctx := context.Background()

	filter := getStackFilterFromOpt(opts.namespace, opts.filter)
	tasks, err := client.TaskList(ctx, types.TaskListOptions{Filters: filter})
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Fprintf(dockerCli.Out(), "Nothing found in stack: %s\n", namespace)
		return nil
	}

	return task.Print(dockerCli, ctx, tasks, idresolver.New(client, opts.noResolve), opts.noTrunc)
}
