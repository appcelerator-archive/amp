package main

import (
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
	"github.com/spf13/cobra"
)

type opts struct {
	version bool
	addr    string
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
			if opts.addr != "" {
				c.SetServer(opts.addr)
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
	cmd.PersistentFlags().StringVarP(&opts.addr, "server", "s", "", "Specify server (host:port)")

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
