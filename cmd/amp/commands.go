package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/cli/command/cluster"
	"github.com/appcelerator/amp/cli/command/completion"
	"github.com/appcelerator/amp/cli/command/config"
	"github.com/appcelerator/amp/cli/command/login"
	"github.com/appcelerator/amp/cli/command/logout"
	"github.com/appcelerator/amp/cli/command/logs"
	"github.com/appcelerator/amp/cli/command/object_store"
	"github.com/appcelerator/amp/cli/command/password"
	"github.com/appcelerator/amp/cli/command/secret"
	"github.com/appcelerator/amp/cli/command/service"
	"github.com/appcelerator/amp/cli/command/settings"
	"github.com/appcelerator/amp/cli/command/stack"
	"github.com/appcelerator/amp/cli/command/team"
	"github.com/appcelerator/amp/cli/command/user"
	"github.com/appcelerator/amp/cli/command/version"
	"github.com/appcelerator/amp/cli/command/whoami"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
)

type opts struct {
	version    bool
	server     string
	skipVerify bool
}

// newRootCommand returns a new instance of the amp cli root command.
func newRootCommand(c cli.Interface) *cobra.Command {
	opts := &opts{}
	cmd := &cobra.Command{
		Use:           "amp",
		Short:         "Deploy, manage, and monitor container stacks.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       "amp version",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.server != "" {
				c.SetServer(opts.server)
			}
			c.SetSkipVerify(opts.skipVerify)

			//print current context
			info(c)

			if cmd.Parent() != nil && cmd.Parent().Use == "cluster" {
				// TODO special case handling for cluster this release
				local := strings.HasPrefix(c.Server(), "127.0.0.1") ||
					strings.HasPrefix(c.Server(), "localhost") ||
					strings.Contains(c.Server(), "local.appcelerator.io")
				switch cmd.Use {
				case "create", "rm", "update":
					if !local {
						return errors.New("cluster operations on remote server are not supported in this release")
					}
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.version {
				showVersion()
				return nil
			}
			cmd.SetOutput(c.Err())
			cmd.HelpFunc()(cmd, args)
			return nil
		},
	}
	cli.SetupRootCommand(cmd)

	cmd.Flags().BoolVarP(&opts.version, "version", "v", false, "Print version information and quit")
	cmd.PersistentFlags().StringVarP(&opts.server, "server", "s", "", "Specify server (host:port)")
	cmd.PersistentFlags().BoolVarP(&opts.skipVerify, "insecure", "k", false, "Control whether amp verifies the server's certificate chain and host name")

	cmd.SetOutput(c.Out())
	addCommands(cmd, c)
	return cmd
}

// addCommands adds the cli commands to the root command that we want to make available for a release.
func addCommands(cmd *cobra.Command, c cli.Interface) {
	cmd.AddCommand(
		// cluster
		cluster.NewClusterCommand(c),

		// completion
		completion.NewCompletionCommand(c),

		// config
		config.NewConfigCommand(c),

		// login
		login.NewLoginCommand(c),

		// logout
		logout.NewLogoutCommand(c),

		// logs
		logs.NewLogsCommand(c),

		// org
		//org.NewOrgCommand(c),

		// password
		password.NewPasswordCommand(c),

		// secret
		secret.NewSecretCommand(c),

		// service
		service.NewServiceCommand(c),

		//settings
		settings.NewSettingsCommand(c),

		// object storage
		object_store.NewObjectStoreCommand(c),

		// stack
		stack.NewStackCommand(c),

		// stats
		//stats.NewStatsCommand(c),

		// team
		team.NewTeamCommand(c),

		// user
		user.NewUserCommand(c),

		// version
		version.NewVersionCommand(c),

		// whoami
		whoami.NewWhoAmICommand(c),
	)
}

func info(c cli.Interface) {
	s := c.Server()
	tkn, err := cli.ReadToken(s)
	if err != nil {
		fmt.Fprintf(c.Err(), "[%s]\n", s)
	}
	if tkn != "" {
		pToken, _ := jwt.ParseWithClaims(tkn, &auth.AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte{}, nil
		})
		if claims, ok := pToken.Claims.(*auth.AuthClaims); ok {
			//if claims.ActiveOrganization != "" {
			//	fmt.Fprintf(c.Err(), "[user %s in organization %s @ %s]\n", claims.AccountName, claims.ActiveOrganization, s)
			//} else {
			fmt.Fprintf(c.Err(), "[user %s @ %s]\n", claims.AccountName, s)
			//}
		}
	}
}
