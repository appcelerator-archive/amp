package main

import (
	"errors"
	"strings"

	//"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/cli/command/cluster"
	"github.com/appcelerator/amp/cli/command/function"
	"github.com/appcelerator/amp/cli/command/login"
	"github.com/appcelerator/amp/cli/command/logout"
	"github.com/appcelerator/amp/cli/command/logs"
	"github.com/appcelerator/amp/cli/command/org"
	"github.com/appcelerator/amp/cli/command/password"
	"github.com/appcelerator/amp/cli/command/service"
	"github.com/appcelerator/amp/cli/command/stack"
	"github.com/appcelerator/amp/cli/command/stats"
	"github.com/appcelerator/amp/cli/command/team"
	"github.com/appcelerator/amp/cli/command/user"
	"github.com/appcelerator/amp/cli/command/version"
	"github.com/appcelerator/amp/cli/command/whoami"
	//"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
)

type opts struct {
	version bool
	server  string
}

// newRootCommand returns a new instance of the amp cli root command.
func newRootCommand(c cli.Interface) *cobra.Command {
	opts := &opts{}
	cmd := &cobra.Command{
		Use:           "amp [OPTIONS] COMMAND [ARG...]",
		Short:         "Deploy, manage, and monitor container stacks and functions.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       "amp version",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.server != "" {
				c.SetServer(opts.server)
			}

			err := info(c)
			if err != nil {
				return err
			}

			if cmd.Parent() != nil && cmd.Parent().Use == "cluster" {
				// TODO special case handling for cluster this release
				local := strings.HasPrefix(c.Server(), "127.0.0.1") ||
					strings.HasPrefix(c.Server(), "localhost")
				if !local {
					return errors.New("only cluster operations with '--server=localhost' supported in this release")
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

	cmd.SetOutput(c.Out())
	addCommands(cmd, c)
	return cmd
}

// addCommands adds the cli commands to the root command that we want to make available for a release.
func addCommands(cmd *cobra.Command, c cli.Interface) {
	cmd.AddCommand(
		// cluster
		cluster.NewClusterCommand(c),

		// function
		function_.NewFunctionCommand(c),

		// login
		login.NewLoginCommand(c),

		// logout
		logout.NewLogoutCommand(c),

		// logs
		logs.NewLogsCommand(c),

		// org
		org.NewOrgCommand(c),

		// password
		password.NewPasswordCommand(c),

		// service
		service.NewServiceCommand(c),

		// stack
		stack.NewStackCommand(c),

		// stats
		stats.NewStatsCommand(c),

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

func info(c cli.Interface) error {
	s := c.Server()
	c.Console().Infof("[%s]\n", s)

	//// TODO: pToken.Claims panics
	//token, err := cli.ReadToken()
	//if err != nil {
	//	c.Console().Infof("[%s] Not logged in. Use `amp login` or `amp user signup`\n.", s)
	//}

	//pToken, _ := jwt.ParseWithClaims(token, &auth.AccountClaims{}, func(t *jwt.Token) (interface{}, error) {
	//	return []byte{}, nil
	//})

	//if claims, ok := pToken.Claims.(*auth.AccountClaims); ok {
	//	if claims.ActiveOrganization != "" {
	//		c.Console().Infof("[%s] user: %s (organization: %s)\n", s, claims.AccountName, claims.ActiveOrganization)
	//	} else {
	//		c.Console().Infof("[%s] user: %s\n", s, claims.AccountName)
	//	}
	//}
	return nil
}
