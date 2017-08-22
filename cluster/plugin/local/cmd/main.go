package main

import (
	"context"
	"fmt"
	"log"
	"os"

	plugin "github.com/appcelerator/amp/cluster/plugin/local"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

const (
	defaultURL     = "unix:///var/run/docker.sock"
	defaultVersion = "1.30"
)

var (
	// Version is set with a linker flag (see Makefile)
	Version string
	// Build is set with a linker flag (see Makefile)
	Build         string
	dockerClient  *client.Client
	defaultLabels = map[string]string{"amp.type.api": "true", "amp.type.route": "true", "amp.type.search": "true", "amp.type.kv": "true", "amp.type.mq": "true", "amp.type.metrics": "true", "amp.type.core": "true", "amp.type.user": "true"}
	opts          = &plugin.RequestOptions{
		Tag:         Version,
		InitRequest: swarm.InitRequest{},
		Labels:      defaultLabels,
		// sane defaults for the local plugin
		Registration:  "none", // overrides current stack default "email"
		Notifications: false,  // just being explicit
	}
)

func initClient(cmd *cobra.Command, args []string) (err error) {
	dockerClient, err = client.NewClient(defaultURL, defaultVersion, nil, nil)
	return
}

func create(cmd *cobra.Command, args []string) {
	// docker swarn init --advertise-addr $interface
	ctx := context.Background()

	if err := plugin.EnsureSwarmExists(ctx, dockerClient, opts); err != nil {
		log.Fatal(err)
	}
	if err := plugin.LabelNode(ctx, dockerClient, opts); err != nil {
		log.Fatal(err)
	}
	if err := plugin.RunAgent(ctx, dockerClient, "install", opts); err != nil {
		log.Fatal(err)
	}
	// use the info command to print json cluster info to stdout
	info(cmd, args)
}

func update(cmd *cobra.Command, args []string) {
	// nothing to do
}

func delete(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	if err := plugin.RunAgent(ctx, dockerClient, "uninstall", opts); err != nil {
		log.Fatal(err)
	}
	if err := plugin.DeleteCluster(ctx, dockerClient, opts); err != nil {
		log.Fatal(err)
	}

	log.Println("cluster deleted")
}

func version(cmd *cobra.Command, args []string) {
	fmt.Printf("Version: %s - Build: %s\n", Version, Build)
}

func info(cmd *cobra.Command, args []string) {
	// docker node inspect self -f '{{.Status.State}}'
	ctx := context.Background()
	coreResp, err := plugin.InfoAMPCore(ctx, dockerClient)
	if err != nil {
		log.Fatal(err)
	}
	userResp, err := plugin.InfoUser(ctx, dockerClient)
	if err != nil {
		log.Fatal(err)
	}
	types, err := dockerClient.Info(ctx)
	if err != nil {
		log.Fatal(err)
	}

	j, err := plugin.InfoToJSON(string(types.Swarm.LocalNodeState), coreResp, userResp)
	if err != nil {
		log.Fatal(err)
	}
	// print json result to stdout
	fmt.Print(j)
}

func main() {
	rootCmd := &cobra.Command{
		Use:               "localplugin",
		Short:             "init/update/destroy an local cluster in Docker swarm mode",
		PersistentPreRunE: initClient,
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "init cluster in swarm mode",
		Run:   create,
	}
	initCmd.PersistentFlags().StringVarP(&opts.InitRequest.ListenAddr, "listen-addr", "l", "0.0.0.0:2377", "Listen address")
	initCmd.PersistentFlags().StringVarP(&opts.InitRequest.AdvertiseAddr, "advertise-addr", "a", "eth0", "Advertise address")
	initCmd.PersistentFlags().BoolVarP(&opts.InitRequest.ForceNewCluster, "force-new-cluster", "", false, "force initialization of a new swarm")
	initCmd.PersistentFlags().BoolVar(&opts.SkipTests, "fast", false, "Skip tests while deploying the core services")
	initCmd.PersistentFlags().BoolVar(&opts.NoMonitoring, "no-monitoring", false, "Don't deploy the monitoring core services")

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "version of the plugin",
		Run:   version,
	}
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "get information about the cluster",
		Run:   info,
	}

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "update the cluster",
		Run:   update,
	}

	destroyCmd := &cobra.Command{
		Use:   "destroy",
		Short: "destroy the cluster",
		Run:   delete,
	}
	destroyCmd.PersistentFlags().BoolVarP(&opts.ForceLeave, "force-leave", "", false, "force leave the swarm")

	// override default value if env vars are set
	val, ok := os.LookupEnv("TAG")
	if ok {
		opts.Tag = val
	}
	val, ok = os.LookupEnv("REGISTRATION")
	if ok {
		opts.Registration = val
	}
	val, ok = os.LookupEnv("NOTIFICATIONS")
	if ok {
		opts.Notifications = val == "true"
	}

	rootCmd.AddCommand(initCmd, versionCmd, infoCmd, updateCmd, destroyCmd)

	_ = rootCmd.Execute()
}
