package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	tablePadding = 3
)

var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string

	// AMP manages the connection and state for the client
	AMP *client.AMP

	// Config is used by command implementations to access the computed client configuration.
	Config                = &client.Configuration{}
	configFile            string
	verbose               bool
	serverAddr            string
	listVersion           = true
	displayConfigFilePath = false

	// RootCmd is the base command for the CLI.
	RootCmd = &cobra.Command{
		Use:   `amp [OPTIONS] COMMAND [arg...]`,
		Short: "Appcelerator Microservice Platform.",
		Long: `Appcelerator Microservice Platform(AMP) is an open-source Container-as-a-Service (CaaS) platform
for managing and monitoring containerized applications and microservices
as part of a unified serverless computing environment.`,
		Run: func(cmd *cobra.Command, args []string) {
			if displayConfigFilePath {
				configFilePath := viper.ConfigFileUsed()
				if configFilePath == "" {
					fmt.Println("No configuration file used (using default configuration)")
				} else {
					fmt.Println(configFilePath)
				}
				cli.Exit(0)
			}
			if listVersion {
				fmt.Printf("amp (cli version: %s, build: %s)\n", Version, Build)
				cli.Exit(0)
			}
			fmt.Println(cmd.UsageString())
		},
	}
)

// All main does is process commands and flags and invoke the app
func main() {
	cobra.OnInitialize(func() {
		cli.InitConfig(configFile, Config, verbose, serverAddr)
		if addr := RootCmd.Flag("server").Value.String(); addr != "" {
			Config.ServerAddress = addr
		}
		if Config.ServerAddress == "" {
			Config.ServerAddress = client.DefaultServerAddress
		}
		if Config.AdminServerAddress == "" {
			Config.AdminServerAddress = client.DefaultAdminServerAddress
		}
		AMP = client.NewAMP(Config, cli.NewLogger(Config.Verbose))
		if !Config.Verbose {
			RootCmd.SilenceErrors = true
			RootCmd.SilenceUsage = true
		}
		cli.AtExit(func() {
			if AMP != nil {
				AMP.Disconnect()
			}
		})
	})

	// infoCmd represents the amp information
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Display AMP version",
		Long:  `Display the current amp version and server information.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("amp (cli version: %s, build: %s)\n", Version, Build)
			fmt.Printf("Server: %s\n", Config.ServerAddress)
		},
	}
	RootCmd.AddCommand(infoCmd)

	// helpCmd displays help about amp Commands
	var helpCmd = &cobra.Command{
		Use:   "help [command]",
		Short: "Help about the command",
		Long:  `Display the help and usage information about AMP commands.`,
		RunE: func(c *cobra.Command, args []string) error {
			cmd, args, e := RootCmd.Find(os.Args[2:])
			if cmd == nil || e != nil || len(args) > 0 {
				return fmt.Errorf("unknown help topic: %v", strings.Join(args, " "))
			}

			helpFunc := cmd.HelpFunc()
			helpFunc(cmd, args)
			return nil
		},
	}
	RootCmd.AddCommand(helpCmd)

	RootCmd.SetUsageTemplate(usageTemplate)
	RootCmd.SetHelpTemplate(helpTemplate)

	RootCmd.PersistentFlags().StringVar(&configFile, "use-config", "", "Specify config file (overrides default at $HOME/.config/amp/amp.yaml)")
	RootCmd.PersistentFlags().BoolVar(&displayConfigFilePath, "config-used", false, "Display config file used (if any)")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	RootCmd.PersistentFlags().StringVar(&serverAddr, "server", "", "Server address")
	RootCmd.PersistentFlags().BoolVarP(&listVersion, "version", "V", false, "Version number")
	RootCmd.PersistentFlags().BoolP("help", "h", false, "Display help")

	cmd, _, err := RootCmd.Find(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		cli.Exit(1)
	}
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		cli.Exit(1)
	}
	cli.Exit(0)
}

var usageTemplate = `Usage:	{{if not .HasSubCommands}}{{.UseLine}}{{end}}{{if .HasSubCommands}}{{ .CommandPath}} COMMAND{{end}}

{{ .Short | trim }}

{{ .Long | trim }}{{if gt .Aliases 0}}

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
