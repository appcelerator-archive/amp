package stack

import (
	"fmt"

	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/filters"
	"github.com/appcelerator/go-docker/cli/cli/command"
	"github.com/appcelerator/go-docker/cli/cli/command/formatter"
	"github.com/appcelerator/go-docker/cli/cli/command/service"
	"github.com/appcelerator/go-docker/cli/opts"
	"golang.org/x/net/context"
)

type ServicesOptions struct {
	Quiet     bool
	Format    string
	Filter    opts.FilterOpt
	Namespace string
}

func RunServices(dockerCli command.Cli, options ServicesOptions) error {
	ctx := context.Background()
	client := dockerCli.Client()

	filter := getStackFilterFromOpt(options.Namespace, options.Filter)
	services, err := client.ServiceList(ctx, types.ServiceListOptions{Filters: filter})
	if err != nil {
		return err
	}

	// if no services in this stack, print message and exit 0
	if len(services) == 0 {
		fmt.Fprintf(dockerCli.Err(), "Nothing found in stack: %s\n", options.Namespace)
		return nil
	}

	info := map[string]formatter.ServiceListInfo{}
	if !options.Quiet {
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

	format := options.Format
	if len(format) == 0 {
		if len(dockerCli.ConfigFile().ServicesFormat) > 0 && !options.Quiet {
			format = dockerCli.ConfigFile().ServicesFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	servicesCtx := formatter.Context{
		Output: dockerCli.Out(),
		Format: formatter.NewServiceListFormat(format, options.Quiet),
	}
	return formatter.ServiceListWrite(servicesCtx, services, info)
}
