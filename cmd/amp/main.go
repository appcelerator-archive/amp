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
		Use:   `amp [OPTIONS] COMMAND [arg...]`,
		Short: "Appcelerator Microservice Platform.",
	}
)

// All main does is process commands and flags and invoke the app
func main() {
	cobra.OnInitialize(func() {
		InitConfig(configFile, &Config, verbose, serverAddr)
		if addr := RootCmd.Flag("server").Value.String(); addr != "" {
			Config.ServerAddress = addr
		}
		if Config.ServerAddress == "" {
			Config.ServerAddress = client.DefaultServerAddress
		}
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
	// infoCmd represents the amp information
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Display amp version and server information",
		Long:  `Display amp version and server information.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("amp (cli version: %s, build: %s)\n", Version, Build)
			fmt.Printf("Server: %s\n", Config.ServerAddress)
		},
	}
	RootCmd.SetUsageTemplate(usageTemplate)
	RootCmd.SetHelpTemplate(helpTemplate)

	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file (default is $HOME/.amp.yaml)")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, `Verbose output`)
	RootCmd.PersistentFlags().StringVar(&serverAddr, "server", "", "Server address")
	RootCmd.AddCommand(configCmd)
	RootCmd.AddCommand(infoCmd)
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		cli.Exit(-1)
	}
	cli.Exit(0)
}

var usageTemplate = `Usage:	{{if not .HasSubCommands}}{{.UseLine}}{{end}}{{if .HasSubCommands}}{{ .CommandPath}} COMMAND{{end}}

{{ .Short | trim }}{{if gt .Aliases 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{ .Example }}{{end}}{{if .HasFlags}}

Options:
{{.Flags.FlagUsages | trimRightSpace}}{{end}}{{ if .HasAvailableSubCommands}}

Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasSubCommands }}

Run '{{.CommandPath}} COMMAND --help' for more information on a command.{{end}}
`

var helpTemplate = `
{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
