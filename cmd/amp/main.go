package main

import (
	"fmt"
	"os"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/spf13/cobra"
)

var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string

	config     client.Configuration
	configFile string
	verbose    bool
)

// All main does is process commands and flags and invoke the app
func main() {
	fmt.Printf("amp (cli version: %s, build: %s)\n", Version, Build)

	cobra.OnInitialize(func() {
		InitConfig(configFile, &config, verbose)
	})

	// createCmd represents the create command
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new AMP swarm",
		Long:  `Create a new AMP swarm for the target environment.`,
		Run: func(cmd *cobra.Command, args []string) {
			a := client.NewAMP(&config)
			a.Create()
		},
	}

	// stopCmd represents the stop command
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a running AMP swarm",
		Long:  `Stop an running AMP swarm.`,
		Run: func(cmd *cobra.Command, args []string) {
			a := client.NewAMP(&config)
			a.Stop()
		},
	}

	// startCmd represents the start command
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start a stopped AMP swarm",
		Long:  `Start a stopped AMP swarm.`,
		Run: func(cmd *cobra.Command, args []string) {
			a := client.NewAMP(&config)
			a.Start()
		},
	}

	// updateCmd represents the update command
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing AMP swarm",
		Long:  `Updated an existing AMP swarm.`,
		Run: func(cmd *cobra.Command, args []string) {
			a := client.NewAMP(&config)
			a.Update()
		},
	}

	// statusCmd represents the status command
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Get status of a running AMP swarm",
		Long:  `Get status of a running AMP swarm.`,
		Run: func(cmd *cobra.Command, args []string) {
			a := client.NewAMP(&config)
			a.Status()
		},
	}

	// configCmd represents the config command
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Display the current configuration",
		Long:  `Display the current configuration, taking into account flags and environment variables.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config)
		},
	}

	// loginCmd represents the login command
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login via github",
		Long:  `Create a github access token and store it in your config file to authenticate further commands`,
		Run: func(cmd *cobra.Command, args []string) {
			a := client.NewAMP(&config)
			cli.Login(a)
		},
	}

	// logsCmd represents the logs command
	logsCmd := &cobra.Command{
		Use:   "logs",
		Short: "Fetch the logs",
		Long:  `Search through all the logs of the system and fetch entries matching provided criteria.`,
		Run: func(cmd *cobra.Command, args []string) {
			amp := client.NewAMP(&config)
			err := cli.Logs(amp, cmd)
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	// TODO logsCmd.Flags().String("timestamp", "", "filter by the given timestamp")
	logsCmd.Flags().String("service_id", "", "filter by the given service id")
	logsCmd.Flags().String("service_name", "", "filter by the given service name")
	logsCmd.Flags().String("message", "", "filter by the given pattern in the message field")
	logsCmd.Flags().String("container_id", "", "filter by the given container id")
	logsCmd.Flags().String("node_id", "", "filter by the given node id")
	logsCmd.Flags().String("from", "-1", "Fetches from the given index")
	logsCmd.Flags().StringP("number", "n", "100", "Number of results")
	logsCmd.Flags().BoolP("short", "s", false, "Displays only the message field")
	logsCmd.Flags().BoolP("follow", "f", false, "Follow log output")

	var statsCmd = &cobra.Command{
		Use:   "stats",
		Short: "Display resource usage statistics",
		Long:  `get statistics on containers, services, nodes about cpu, memory, io, net.`,
		Run: func(cmd *cobra.Command, args []string) {
			amp := client.NewAMP(&config)
			err := cli.Stats(amp, cmd, args)
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	statsCmd.Flags().Bool("container", false, "display stats on containers")
	statsCmd.Flags().Bool("service", false, "displat stats on services")
	statsCmd.Flags().Bool("node", false, "display stats on nodes")
	statsCmd.Flags().Bool("task", false, "display stats on tasks")
	//metrics
	statsCmd.Flags().Bool("cpu", false, "display cpu stats")
	//statsCmd.Flags().Bool("mem", false, "display memory stats")
	//statsCmd.Flags().Bool("io", false, "display memory stats")
	//statsCmd.Flags().Bool("net", false, "display memory stats")
	//historic
	statsCmd.Flags().String("period", "", "historic period of metrics extraction, duration + time unit")
	statsCmd.Flags().String("since", "", "date defining when begin the historic metrics extraction, format: YYYY-MM-DD HH:MM:SS.mmm")
	statsCmd.Flags().String("until", "", "date defining when stop the historic metrics extraction, format: YYYY-MM-DD HH:MM:SS.mmm")
	//statsCmd.Flags().String("time-unit", "", "historic extraction group can be: s:seconds, m:minutes, h:hours, d:days, w:weeks")
	//filters:
	statsCmd.Flags().String("container-id", "", "filter on container id")
	statsCmd.Flags().String("container-name", "", "filter on container name")
	statsCmd.Flags().String("image", "", "filter on container image name")
	statsCmd.Flags().String("service-name", "", "filter on service name")
	statsCmd.Flags().String("service-id", "", "filter on service id")
	statsCmd.Flags().String("task-name", "", "filter on task name")
	statsCmd.Flags().String("task-id", "", "filter on task id")
	statsCmd.Flags().String("datacenter", "", "filter on datacenter")
	statsCmd.Flags().String("host", "", "filter on host")
	statsCmd.Flags().String("node-id", "", "filter on node id")

	// This represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "amp",
		Short: "Manage an AMP swarm",
		Long:  `Manage an AMP swarm.`,
	}
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.amp.yaml)")
	rootCmd.PersistentFlags().String("target", "local", `target environment ("local"|"virtualbox"|"aws")`)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, `verbose output`)

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(statsCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
