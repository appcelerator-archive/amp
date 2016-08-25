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

	// Config is used by command implementations to access the current client configuration.
	Config     client.Configuration
	configFile string
	verbose    bool

	// RootCmd is the base command for the CLI.
	RootCmd = &cobra.Command{
		Use:   "amp",
		Short: "Manage an AMP swarm",
		Long:  `Manage an AMP swarm.`,
	}
)

// All main does is process commands and flags and invoke the app
func main() {
	fmt.Printf("amp (cli version: %s, build: %s)\n", Version, Build)

	cobra.OnInitialize(func() {
		InitConfig(configFile, &Config, verbose)
	})

	// createCmd represents the create command
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new AMP swarm",
		Long:  `Create a new AMP swarm for the target environment.`,
		Run: func(cmd *cobra.Command, args []string) {
			a := client.NewAMP(&Config)
			a.Create()
		},
	}

	// stopCmd represents the stop command
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a running AMP swarm",
		Long:  `Stop an running AMP swarm.`,
		Run: func(cmd *cobra.Command, args []string) {
			a := client.NewAMP(&Config)
			a.Stop()
		},
	}

	// startCmd represents the start command
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start a stopped AMP swarm",
		Long:  `Start a stopped AMP swarm.`,
		Run: func(cmd *cobra.Command, args []string) {
			a := client.NewAMP(&Config)
			a.Start()
		},
	}

	// updateCmd represents the update command
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing AMP swarm",
		Long:  `Updated an existing AMP swarm.`,
		Run: func(cmd *cobra.Command, args []string) {
			a := client.NewAMP(&Config)
			a.Update()
		},
	}

	// statusCmd represents the status command
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Get status of a running AMP swarm",
		Long:  `Get status of a running AMP swarm.`,
		Run: func(cmd *cobra.Command, args []string) {
			a := client.NewAMP(&Config)
			a.Status()
		},
	}

	// configCmd represents the Config command
	configCmd := &cobra.Command{
		Use:   "Config",
		Short: "Display the current configuration",
		Long:  `Display the current configuration, taking into account flags and environment variables.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Config)
		},
	}

	// loginCmd represents the login command
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login via github",
		Long:  `Create a github access token and store it in your Config file to authenticate further commands`,
		Run: func(cmd *cobra.Command, args []string) {
			a := client.NewAMP(&Config)
			cli.Login(a)
		},
	}

	// logsCmd represents the logs command
	logsCmd := &cobra.Command{
		Use:   "logs",
		Short: "Fetch the logs",
		Long:  `Search through all the logs of the system and fetch entries matching provided criteria.`,
		Run: func(cmd *cobra.Command, args []string) {
			amp := client.NewAMP(&Config)
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
	logsCmd.Flags().String("n", "100", "Number of results")
	logsCmd.Flags().Bool("short", false, "Displays only the message field")
	logsCmd.Flags().Bool("f", false, "Follow log output")

	// This represents the base command when called without any subcommands
	RootCmd.PersistentFlags().StringVar(&configFile, "Config", "", "Config file (default is $HOME/.amp.yaml)")
	RootCmd.PersistentFlags().String("target", "local", `target environment ("local"|"virtualbox"|"aws")`)
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, `verbose output`)

	RootCmd.AddCommand(createCmd)
	RootCmd.AddCommand(stopCmd)
	RootCmd.AddCommand(startCmd)
	RootCmd.AddCommand(updateCmd)
	RootCmd.AddCommand(statusCmd)
	RootCmd.AddCommand(configCmd)
	RootCmd.AddCommand(loginCmd)
	RootCmd.AddCommand(logsCmd)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
