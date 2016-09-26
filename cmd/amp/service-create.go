package main

import (
	"fmt"
	"strings"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
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

	// service labels
	serviceLabels []string

	// container labels
	containerLabels []string

	// ports
	publishSpecs []string
)

func init() {
	flags := createCmd.Flags()
	flags.StringVar(&name, "name", name, "Service name")
	flags.StringSliceVarP(&publishSpecs, "publish", "p", publishSpecs, "Publish a service externally. Format: [published-name|published-port:]internal-service-port[/protocol], i.e. '80:3000/tcp' or 'admin:3000'")
	flags.Uint64Var(&replicas, "replicas", replicas, "Number of tasks (default none)")
	flags.StringSliceVarP(&env, "env", "e", env, "Set environment variables (default [])")
	flags.StringSliceVarP(&serviceLabels, "label", "l", serviceLabels, "Set service labels (default [])")
	flags.StringSliceVar(&containerLabels, "container-label", containerLabels, "Set container labels for service replicas (default [])")

	ServiceCmd.AddCommand(createCmd)
}

func create(amp *client.AMP, cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		// TODO use standard errors and print usage
		return fmt.Errorf("\"amp service create\" requires at least 1 argument(s)")
	}

	image = args[0]

	parsedSpecs, err := parsePublishSpecs(publishSpecs)
	if err != nil {
		return err
	}

	spec := &service.ServiceSpec{
		Image:           image,
		Name:            name,
		Replicas:        replicas,
		Env:             env,
		Labels:          stringmap(serviceLabels),
		ContainerLabels: stringmap(containerLabels),
		PublishSpecs:    parsedSpecs,
	}

	request := &service.ServiceCreateRequest{
		ServiceSpec: spec,
	}

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

func parsePublishSpecs(specs []string) ([]*service.PublishSpec, error) {
	publishSpecs := []*service.PublishSpec{}
	for _, input := range specs {
		publishSpec, err := service.ParsePublishSpec(input)
		if err != nil {
			return nil, err
		}
		publishSpecs = append(publishSpecs, &publishSpec)

	}
	return publishSpecs, nil
}
