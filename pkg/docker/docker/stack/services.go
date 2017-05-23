package stack

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/cli/command"
	"github.com/docker/docker/cli/command/formatter"
	"github.com/docker/docker/cli/command/service"
	"github.com/docker/docker/opts"
	"golang.org/x/net/context"
)

type servicesOptions struct {
	quiet     bool
	format    string
	filter    opts.FilterOpt
	namespace string
}

// NewServicesOptions {AMP}: add constructor for this private struct
func NewServicesOptions(quiet bool, format string, filter opts.FilterOpt, namespace string) servicesOptions {
	return servicesOptions{
		quiet:     quiet,
		format:    format,
		filter:    filter,
		namespace: namespace,
	}
}

// RunServices {amp}: make it public
func RunServices(dockerCli *command.DockerCli, opts servicesOptions) error {
	ctx := context.Background()
	client := dockerCli.Client()

	filter := getStackFilterFromOpt(opts.namespace, opts.filter)
	services, err := client.ServiceList(ctx, types.ServiceListOptions{Filters: filter})
	if err != nil {
		return err
	}

	out := dockerCli.Out()

	// if no services in this stack, print message and exit 0
	if len(services) == 0 {
		fmt.Fprintf(out, "Nothing found in stack: %s\n", opts.namespace)
		return nil
	}

	info := map[string]formatter.ServiceListInfo{}
	if !opts.quiet {
		taskFilter := filters.NewArgs()
		for _, service := range services {
			taskFilter.Add("service", service.ID)
		}

		tasks, err := client.TaskList(ctx, types.TaskListOptions{Filters: taskFilter})
		if err != nil {
			return err
		}

		nodes, err := client.NodeList(ctx, types.NodeListOptions{})
		if err != nil {
			return err
		}

		info = service.GetServicesStatus(services, nodes, tasks)
	}

	format := opts.format
	if len(format) == 0 {
		if len(dockerCli.ConfigFile().ServicesFormat) > 0 && !opts.quiet {
			format = dockerCli.ConfigFile().ServicesFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	servicesCtx := formatter.Context{
		Output: dockerCli.Out(),
		Format: formatter.NewServiceListFormat(format, opts.quiet),
	}
	return formatter.ServiceListWrite(servicesCtx, services, info)
}
