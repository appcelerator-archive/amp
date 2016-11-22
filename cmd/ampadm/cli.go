package main

import (
	"github.com/spf13/cobra"
	"os"
)

const (
	padding = 3
)

var (
	serverAddr = ""
	configFile = ""
)

func (g *clusterClient) cli() {
	if err := client.init(); err != nil {
		g.printfc(colError, "Init error: %v\n", err)
		os.Exit(1)
	}
	RootCmd.PersistentFlags().BoolVarP(&g.verbose, "verbose", "v", false, `Verbose output`)
	RootCmd.PersistentFlags().BoolVar(&g.silence, "silence", false, `No output`)
	RootCmd.PersistentFlags().BoolVar(&g.debug, "debug", false, `display debug information`)
	RootCmd.PersistentFlags().StringVar(&serverAddr, "server", "127.0.0.1:31315", `define the remote server address`)
	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file (default is $HOME/.amp.yaml)")
	cobra.OnInitialize(func() {
		g.initConfiguration(configFile, serverAddr)
	})

	// versionCmd displays the agreed version
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display the version number of adm-server",
		Long:  `Display the version number of adm-server`,
		Run: func(cmd *cobra.Command, args []string) {
			g.printfc(colSuccess, "ampadm version: %s, build: %s)\n", Version, Build)
		},
	}
	RootCmd.AddCommand(versionCmd)

	//Execute command
	cmd, _, err := RootCmd.Find(os.Args[1:])
	if err != nil {
		g.fatalc("Arg error: %v\n", err)
	}
	if err := cmd.Execute(); err != nil {
		g.fatalc("Error during: %s: %v\n", cmd.Name(), err)
	}

	os.Exit(0)
}
