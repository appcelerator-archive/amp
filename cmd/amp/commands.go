package main

import (
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/cli/command/cluster"
	"github.com/appcelerator/amp/cli/command/login"
	"github.com/appcelerator/amp/cli/command/logout"
	"github.com/appcelerator/amp/cli/command/logs"
	"github.com/appcelerator/amp/cli/command/org"
	"github.com/appcelerator/amp/cli/command/password"
	"github.com/appcelerator/amp/cli/command/stats"
	"github.com/appcelerator/amp/cli/command/team"
	"github.com/appcelerator/amp/cli/command/user"
	"github.com/appcelerator/amp/cli/command/version"
	"github.com/appcelerator/amp/cli/command/whoami"
	"github.com/spf13/cobra"
)

type opts struct {
	version bool
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
		Run: func(cmd *cobra.Command, args []string) {
			if opts.version {
				showVersion()
				return
			}
			cmd.SetOutput(c.Err())
			cmd.HelpFunc()(cmd, args)
		},
	}
	cli.SetupRootCommand(cmd)

	flags := cmd.Flags()
	flags.BoolVarP(&opts.version, "version", "v", false, "Print version information and quit")

	cmd.SetOutput(c.Out())
	addCommands(cmd, c)
	return cmd
}

// addCommands adds the cli commands to the root command that we want to make available for a release.
func addCommands(cmd *cobra.Command, c cli.Interface) {
	cmd.AddCommand(
		// bootstrap
		cluster.NewClusterCommand(c),

		// login
		login.NewLoginCommand(c),

		// logs
		logs.NewLogsCommand(c),

		//logout
		logout.NewLogoutCommand(c),

		// org
		org.NewOrgCommand(c),

		//password
		password.NewPasswordCommand(c),

		//team
		team.NewTeamCommand(c),

		// user
		user.NewUserCommand(c),

		// Stats
		stats.NewStatsCommand(c),

		// version
		version.NewVersionCommand(c),

		//whoami
		whoami.NewWhoAmICommand(c),
	)
}
