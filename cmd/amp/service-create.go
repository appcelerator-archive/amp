package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:   "create [OPTIONS] IMAGE [CMD] [ARG...]",
		Short: "Create a service",
		Long:  `Create a new service`,
		Run: func(cmd *cobra.Command, args []string) {
			err := create(AMP, cmd, args)
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	// service image
	image string

	// service name
	name string

	// number of tasks
	replicas uint64 = 1

	// environment variables
	env []string

	// ports
	publishPorts []string
	exposePorts  []string
)

func init() {
	flags := createCmd.Flags()
	flags.StringVar(&name, "name", name, "Service name")
	flags.Uint64Var(&replicas, "replicas", replicas, "Number of tasks (default none)")
	flags.StringSliceVarP(&env, "env", "e", env, "Set environment variables (default [])")
	flags.StringSliceVarP(&publishPorts, "publish", "p", publishPorts, "Publish a container's port to the host. Format: [hostPort:]containerPort[/protocol], i.e. '80:80/tcp'")
	flags.StringSliceVar(&exposePorts, "expose", exposePorts, "Publish a container's port to the host. Format: [hostPort:]containerPort[/protocol], i.e. '80:80/tcp'")

	ServiceCmd.AddCommand(createCmd)
}

func create(amp *client.AMP, cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		// TODO use standard errors and print usage
		return fmt.Errorf("\"amp service create\" requires at least 1 argument(s)")
	}

	image = args[0]
	fmt.Println(args)
	fmt.Println(stringify(cmd))

	spec := &service.ServiceSpec{
		Image:        image,
		Name:         name,
		Replicas:     replicas,
		Env:          stringmap(env),
		PublishPorts: parsePublishPorts(publishPorts),
		ExposePorts:  parseExposePorts(exposePorts),
	}

	request := &service.ServiceCreateRequest{
		ServiceSpec: spec,
	}

	fmt.Println(request)

	client := service.NewServiceClient(amp.Conn)
	reply, err := client.Create(context.Background(), request)
	if err != nil {
		return err
	}

	fmt.Println(reply)

	return nil
}

func stringmap(a []string) map[string]string {
	m := make(map[string]string)
	for _, e := range a {
		parts := strings.Split(e, "=")
		m[parts[0]] = parts[1]
	}
	return m
}

func stringify(cmd *cobra.Command) string {
	return fmt.Sprintf("{ name: %s, replicas: %d, env: %v }",
		name, replicas, env)
}

func parsePublishPorts(ports []string) *service.Ports {
	return &service.Ports{}
}

func parseExposePorts(ports []string) *service.Ports {
	return &service.Ports{}
}
