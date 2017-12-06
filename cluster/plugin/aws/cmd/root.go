package cmd

import (
	"github.com/appcelerator/amp/cluster/plugin/aws/plugin"
	"github.com/spf13/cobra"
)

// NewRootCommand returns a new instance of the root command
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "awsplugin",
		Short: "Manage AWS cluster in Docker swarm mode",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Instantiate a new AWSPlugin instance
			plugin.AWS = plugin.New(plugin.Config)
		},
	}

	// Flags
	cmd.PersistentFlags().StringVarP(&plugin.Config.Region, "region", "r", "", "aws region")
	cmd.PersistentFlags().StringVarP(&plugin.Config.StackName, "stackname", "n", "", "aws stack name")
	cmd.PersistentFlags().StringSliceVarP(&plugin.Config.Params, "parameter", "p", []string{}, "parameter")
	cmd.PersistentFlags().BoolVarP(&plugin.Config.Sync, "sync", "s", true, "block until operation is complete")
	cmd.PersistentFlags().StringVar(&plugin.Config.AccessKeyId, "access-key-id", "", "access key id (for example, AKIAIOSFODNN7EXAMPLE)")
	cmd.PersistentFlags().StringVar(&plugin.Config.SecretAccessKey, "secret-access-key", "", "secret access key (for example, wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY)")
	cmd.PersistentFlags().StringVar(&plugin.Config.Profile, "profile", "default", "credential profile")

	// Sub commands
	cmd.AddCommand(NewCreateCommand())
	cmd.AddCommand(NewRemoveCommand())
	cmd.AddCommand(NewUpdateCommand())
	return cmd
}
