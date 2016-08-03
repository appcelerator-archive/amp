package main

import (
	"fmt"
	"os"

	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string

	config     cli.Config
	configFile string
	verbose    bool
)

// All main does is process commands and flags and invoke the app
func main() {
	fmt.Printf("amp (cli version: %s, build: %s)\n", Version, Build)

	cobra.OnInitialize(func() {
		cli.InitConfig(configFile, &config, verbose)
	})

	// createCmd represents the create command
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new AMP swarm",
		Long:  `Create a new AMP swarm for the target environment.`,
		Run: func(cmd *cobra.Command, args []string) {
			a := cli.NewAMP(&config)
			a.Create()
		},
	}

	// stopCmd represents the stop command
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a running AMP swarm",
		Long:  `Stop an running AMP swarm.`,
		Run: func(cmd *cobra.Command, args []string) {
			a := cli.NewAMP(&config)
			a.Stop()
		},
	}

	// startCmd represents the start command
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start a stopped AMP swarm",
		Long:  `Start a stopped AMP swarm.`,
		Run: func(cmd *cobra.Command, args []string) {
			a := cli.NewAMP(&config)
			a.Start()
		},
	}

	// updateCmd represents the update command
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing AMP swarm",
		Long:  `Updated an existing AMP swarm.`,
		Run: func(cmd *cobra.Command, args []string) {
			a := cli.NewAMP(&config)
			a.Update()
		},
	}

	// statusCmd represents the status command
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Get status of a running AMP swarm",
		Long:  `Get status of a running AMP swarm.`,
		Run: func(cmd *cobra.Command, args []string) {
			a := cli.NewAMP(&config)
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

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
