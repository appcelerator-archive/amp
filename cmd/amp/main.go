package main

import (
	"fmt"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/spf13/cobra"
)

var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string

	// AMP manages the connection and state for the client
	AMP *client.AMP

	// Config is used by command implementations to access the computed client configuration.
	Config     client.Configuration
	configFile string
	verbose    bool
	serverAddr string

	// RootCmd is the base command for the CLI.
	RootCmd = &cobra.Command{
		Use:   "amp",
		Short: "AMP CLI",
		Long:  `AMP CLI.`,
	}
)

// All main does is process commands and flags and invoke the app
func main() {
	fmt.Printf("amp (cli version: %s, build: %s)\n", Version, Build)

	cobra.OnInitialize(func() {
		InitConfig(configFile, &Config, verbose, serverAddr)
		fmt.Println("Server: " + Config.ServerAddress)
		AMP = client.NewAMP(&Config)
		AMP.Connect()
		cli.AtExit(func() {
			if AMP != nil {
				AMP.Disconnect()
			}
		})
	})

	// configCmd represents the Config command
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Display the current configuration",
		Long:  `Display the current configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Config)
		},
	}

	RootCmd.PersistentFlags().StringVar(&configFile, "Config", "", "Config file (default is $HOME/.amp.yaml)")
	RootCmd.PersistentFlags().String("target", "local", `target environment ("local"|"virtualbox"|"aws")`)
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, `verbose output`)
	RootCmd.PersistentFlags().StringVar(&serverAddr, "server", client.DefaultServerAddress, "Server address")

	RootCmd.AddCommand(configCmd)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		cli.Exit(-1)
	}
	cli.Exit(0)
}
