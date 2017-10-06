package stack

import (
	"fmt"

	"golang.org/x/net/context"
	"github.com/appcelerator/amp/docker/cli/opts"
	"docker.io/go-docker/api/types/filters"
	"github.com/appcelerator/amp/docker/cli/cli/command"
	"docker.io/go-docker/api/types"
	"github.com/appcelerator/amp/docker/cli/cli/command/formatter"
	"github.com/appcelerator/amp/docker/cli/cli/command/service"
)

type servicesOptions struct {
	quiet     bool
	format    string
	filter    opts.FilterOpt
	namespace string
}

func NewServicesOptions(quiet bool, format string, filter opts.FilterOpt, namespace string) servicesOptions {
	return servicesOptions{
		quiet:     quiet,
		format:    format,
		filter:    filter,
		namespace: namespace,
	}
}

func RunServices(dockerCli command.Cli, options servicesOptions) error {
	ctx := context.Background()
	client := dockerCli.Client()

	filter := getStackFilterFromOpt(options.namespace, options.filter)
	services, err := client.ServiceList(ctx, types.ServiceListOptions{Filters: filter})
	if err != nil {
		return err
	}

	// if no services in this stack, print message and exit 0
	if len(services) == 0 {
		fmt.Fprintf(dockerCli.Err(), "Nothing found in stack: %s\n", options.namespace)
		return nil
	}

	info := map[string]formatter.ServiceListInfo{}
	if !options.quiet {
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

	format := options.format
	if len(format) == 0 {
		if len(dockerCli.ConfigFile().ServicesFormat) > 0 && !options.quiet {
			format = dockerCli.ConfigFile().ServicesFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	servicesCtx := formatter.Context{
		Output: dockerCli.Out(),
		Format: formatter.NewServiceListFormat(format, options.quiet),
	}
	return formatter.ServiceListWrite(servicesCtx, services, info)
}
