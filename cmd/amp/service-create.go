package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/docker/docker/cli/command"
	cliflags "github.com/docker/docker/cli/flags"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	serviceCreateCmd = &cobra.Command{
		Use:   "create [OPTIONS] IMAGE [CMD] [ARG...]",
		Short: "Create a new service",
		Long:  `Create a new service.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return serviceCreate(AMP, cmd, args)
		},
	}

	// service image
	image string

	// service name
	name string

	// service mode
	mode = "replicated"

	// number of tasks
	replicas uint64

	// environment variables
	env = []string{}

	// service labels
	serviceLabels = []string{}

	// container labels
	containerLabels = []string{}

	// ports
	publishSpecs = []string{}

	// network
	networks = []string{}

	// send registry auth
	registryAuth = false
)

func init() {
	flags := serviceCreateCmd.Flags()
	flags.StringVar(&name, "name", name, "Service name")
	flags.StringSliceVarP(&publishSpecs, "publish", "p", publishSpecs, "Publish a service externally. Format: [published-name|published-port:]internal-service-port[/protocol], ex: '80:3000/tcp' or 'admin:3000'")
	flags.StringVar(&mode, "mode", mode, "Service mode (replicated or global)")
	flags.Uint64Var(&replicas, "replicas", replicas, "Number of tasks")
	flags.StringSliceVarP(&env, "env", "e", env, "Set environment variables (default [])")
	flags.StringSliceVarP(&serviceLabels, "label", "l", serviceLabels, "Set service labels (default [])")
	flags.StringSliceVar(&containerLabels, "container-label", containerLabels, "Set container labels for service replicas (default [])")
	flags.StringSliceVar(&networks, "network", networks, "Set service networks attachment (default [])")
	flags.BoolVar(&registryAuth, "with-registry-auth", false, "Send registry authentication details to Swarm agents")

	ServiceCmd.AddCommand(serviceCreateCmd)
}

func serviceCreate(amp *client.AMP, cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		// TODO use standard errors and print usage
		log.Fatal("\"amp service create\" requires at least 1 argument(s)")
	}

	image = args[0]

	parsedSpecs, err := parsePublishSpecs(publishSpecs)
	if err != nil {
		return err
	}

	parsedNetworks, err := parseNetworks(networks)
	if err != nil {
		return err
	}

	// add service mode to spec
	var swarmMode service.SwarmMode
	switch mode {
	case "replicated":
		if replicas < 1 {
			// if replicated then must have at least 1 replica
			replicas = 1
		}
		swarmMode = &service.ServiceSpec_Replicated{
			Replicated: &service.ReplicatedService{Replicas: replicas},
		}
	case "global":
		if replicas != 0 {
			// global mode can't specify replicas (only allowed 1 per node)
			log.Fatal("Replicas can only be used with replicated mode")
		}
		swarmMode = &service.ServiceSpec_Global{
			Global: &service.GlobalService{},
		}
	default:
		log.Fatalf("Invalid option for mode: %s", mode)
	}

	spec := &service.ServiceSpec{
		Image:           image,
		Name:            name,
		Env:             env,
		Mode:            swarmMode,
		Labels:          stringmap(serviceLabels),
		ContainerLabels: stringmap(containerLabels),
		PublishSpecs:    parsedSpecs,
		Networks:        parsedNetworks,
	}

	request := &service.ServiceCreateRequest{
		ServiceSpec: spec,
	}

	ctx := context.Background()

	// only send auth if flag was set
	if registryAuth {
		dockerCli := command.NewDockerCli(os.Stdin, os.Stdout, os.Stderr)
		opts := cliflags.NewClientOptions()
		err := dockerCli.Initialize(opts)
		if err != nil {
			return err
		}
		// Retrieve encoded auth token from the image reference
		encodedAuth, err := command.RetrieveAuthTokenFromImage(ctx, dockerCli, image)
		if err != nil {
			return err
		}
		spec.RegistryAuth = encodedAuth
	}

	client := service.NewServiceClient(amp.Conn)
	reply, err := client.Create(ctx, request)
	if err != nil {
		return err
	}

	fmt.Println(reply.Id)
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

func parseNetworks(specs []string) ([]*service.NetworkAttachment, error) {
	networks := []*service.NetworkAttachment{}
	for _, input := range specs {
		network := service.ParseNetwork(input)
		networks = append(networks, network)

	}
	return networks, nil
}
