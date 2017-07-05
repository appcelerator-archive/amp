package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ampctl",
		Short: "Run commands in target amp cluster",
	}

	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Run validation tests on the cluster",
		Run:   checks,
	}
	checkCmd.Flags().BoolVar(&checksOpts.version, "version", false, "check Docker version")
	checkCmd.Flags().BoolVar(&checksOpts.scheduling, "scheduling", false, "check Docker service scheduling")
	checkCmd.Flags().BoolVarP(&checksOpts.all, "all", "a", false, "all tests")

	rootCmd.AddCommand(checkCmd)

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Set up amp services in swarm environment",
		Run:   install,
	}
	rootCmd.AddCommand(installCmd)

	monitorCmd := &cobra.Command{
		Use:   "monitor",
		Short: "Monitor swarm events",
		Run:   monitor,
	}
	rootCmd.AddCommand(monitorCmd)

	_ = rootCmd.Execute()
}
