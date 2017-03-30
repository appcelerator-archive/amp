package main

import (
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/cli/command/version"
	"github.com/spf13/cobra"
)

type Opts struct {
	version bool
}

// newRootCommand returns a new instance of the amp cli root command.
func newRootCommand(c cli.Interface) *cobra.Command {
	opts := &Opts{}

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
		// version
		version.NewVersionCommand(c),
	)
}
