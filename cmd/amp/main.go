package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/appcelerator/amp/api/authn"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
	AMP *cli.AMP

	mgr *cli.CmdManager

	// Config is used by command implementations to access the computed client configuration.
	Config                = &cli.Configuration{}
	configFile            string
	verbose               bool
	serverAddr            string
	listVersion           = true
	displayConfigFilePath = false

	// RootCmd is the base command for the CLI.
	RootCmd = &cobra.Command{
		Use:     "amp",
		Short:   "Appcelerator Microservice Platform",
		Example: "amp org \namp kv get foo",
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

	helpCmd = &cobra.Command{
		Use:     "help",
		Short:   "Help about the command",
		Example: " ",
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

	infoCmd = &cobra.Command{
		Use:     "info",
		Short:   "Display AMP version",
		Example: " ",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("amp (cli version: %s, build: %s)\n", Version, Build)
			fmt.Printf("Server: %s\n", Config.AmpAddress)
		},
	}

	loginCmd = &cobra.Command{
		Use:     "login",
		Short:   "Login to account",
		Example: "--name=jdoe --password=p@s5wrd",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return AMP.Connect()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return login(AMP, cmd)
		},
	}

	switchCmd = &cobra.Command{
		Use:     "switch",
		Short:   "Switch account",
		Example: "--name=swatkats",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return AMP.Connect()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return switchAccount(AMP, cmd)
		},
	}

	whoAmICmd = &cobra.Command{
		Use:     "whoami",
		Short:   "Display currently logged-in user",
		Example: " ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return whoAmI()
		},
	}

	logoutCmd = &cobra.Command{
		Use:     "logout",
		Short:   "Logout current user",
		Example: " ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return logout()
		},
	}

	username string
	password string
)

func init() {
	RootCmd.AddCommand(infoCmd)

	RootCmd.AddCommand(helpCmd)
	RootCmd.AddCommand(loginCmd)
	RootCmd.AddCommand(switchCmd)
	RootCmd.AddCommand(whoAmICmd)
	RootCmd.AddCommand(logoutCmd)

	RootCmd.SetUsageTemplate(usageTemplate)
	RootCmd.SetHelpTemplate(helpTemplate)

	RootCmd.PersistentFlags().StringVar(&configFile, "use-config", "", "Specify config file (overrides default at $HOME/.config/amp/amp.yaml)")
	RootCmd.PersistentFlags().BoolVar(&displayConfigFilePath, "config-used", false, "Display config file used (if any)")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	RootCmd.PersistentFlags().StringVar(&serverAddr, "server", "", "Server address")
	RootCmd.PersistentFlags().BoolVarP(&listVersion, "version", "V", false, "Version number")
	RootCmd.PersistentFlags().BoolP("help", "h", false, "Display help")

	loginCmd.Flags().StringVar(&username, "name", "", "Account Name")
	loginCmd.Flags().StringVar(&password, "password", "", "Password")

	switchCmd.Flags().StringVar(&username, "name", username, "Account Name")
}

// All main does is process commands and flags and invoke the app
func main() {
	cobra.OnInitialize(func() {
		cli.InitConfig(configFile, Config, verbose, serverAddr)
		if addr := RootCmd.Flag("server").Value.String(); addr != "" {
			Config.AmpAddress = addr
		}
		AMP = cli.NewAMP(Config, cli.NewLogger(Config.Verbose))
		if !Config.Verbose {
			RootCmd.SilenceErrors = true
			RootCmd.SilenceUsage = true
		}

		//initialize command manager
		mgr = cli.NewCmdManager("")

		cli.AtExit(func() {
			if AMP != nil {
				AMP.Disconnect()
			}
		})
	})

	cmd, _, err := RootCmd.Find(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		cli.Exit(1)
	}
	if err := cmd.Execute(); err != nil {
		//fmt.Println(err)
		//cli.Exit(1)
		mgr.Fatal(grpc.ErrorDesc(err))
	}
	cli.Exit(0)
}

// login validates the input command line arguments and allows login to an existing account
// by invoking the corresponding rpc/storage method
func login(amp *cli.AMP, cmd *cobra.Command) error {
	if cmd.Flag("name").Changed {
		username = cmd.Flag("name").Value.String()
	} else {
		fmt.Print("username: ")
		username = getName()
	}
	if cmd.Flag("password").Changed {
		password = cmd.Flag("password").Value.String()
	} else {
		password = getPassword()
	}

	request := &account.LogInRequest{
		Name:     username,
		Password: password,
	}
	accClient := account.NewAccountClient(amp.Conn)
	header := metadata.MD{}
	_, err := accClient.Login(context.Background(), request, grpc.Header(&header))
	if err != nil {
		mgr.Fatal(grpc.ErrorDesc(err))
	}
	if err := cli.SaveToken(header); err != nil {
		mgr.Fatal(grpc.ErrorDesc(err))
	}
	mgr.Success("Welcome back, %s!", username)
	return nil
}

// switchAccount validates the input command line arguments and switches from personal account to an organization account
// by invoking the corresponding rpc/storage method
func switchAccount(amp *cli.AMP, cmd *cobra.Command) error {
	if cmd.Flag("name").Changed {
		username = cmd.Flag("name").Value.String()
	} else {
		fmt.Print("account: ")
		username = getName()
	}

	request := &account.SwitchRequest{
		Account: username,
	}
	accClient := account.NewAccountClient(amp.Conn)
	header := metadata.MD{}
	_, err := accClient.Switch(context.Background(), request, grpc.Header(&header))
	if err != nil {
		//manager.fatalf(grpc.ErrorDesc(err))
		//return
		mgr.Fatal(grpc.ErrorDesc(err))
	}
	if err := cli.SaveToken(header); err != nil {
		//return err
		mgr.Fatal(grpc.ErrorDesc(err))
	}
	mgr.Success("Your are now logged in as: %s", username)
	return nil
}

// whoAmI validates the input command line arguments and displays the current account
// by invoking the corresponding rpc/storage method
func whoAmI() error {
	token, err := cli.ReadToken()
	if err != nil {
		//manager.fatalf("You are not logged in.")
		//return
		mgr.Fatal("you are not logged in")
	}
	pToken, _ := jwt.ParseWithClaims(token, &authn.AccountClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte{}, nil
	})
	if claims, ok := pToken.Claims.(*authn.AccountClaims); ok {
		if claims.ActiveOrganization != "" {
			mgr.Success("Logged in as organization %s (on behalf of user %s).", claims.ActiveOrganization, claims.AccountName)
		} else {
			mgr.Success("Logged in as user %s.", claims.AccountName)
		}
	}
	return nil
}

// logout validates the input command line arguments and logs out of the current account
// by invoking the corresponding rpc/storage method
func logout() error {
	err := cli.RemoveToken()
	if err != nil {
		mgr.Fatal(grpc.ErrorDesc(err))
	}
	mgr.Success("You have been successfully logged out!")
	return nil
}

var usageTemplate = `Usage: {{if not .HasSubCommands}}{{.UseLine}}{{end}}{{if .HasSubCommands}}{{ .CommandPath}} COMMAND{{end}}

{{ .Short | trim }}{{if gt .Aliases 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{if not .HasSubCommands}}{{.UseLine}}{{end}} {{ .Example }}{{end}}{{if .HasFlags}}

Options:
{{.Flags.FlagUsages | trimRightSpace}}{{end}}{{ if .HasAvailableSubCommands}}

Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasSubCommands }}

Run '{{.CommandPath}} COMMAND --help' for more information on a command.{{end}}
`

var helpTemplate = `
{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
