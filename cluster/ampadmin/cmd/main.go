package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "envoy",
		Short: "run commands in target cluster",
	}

	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "run validation tests on the cluster",
		Run:   checks,
	}
	checkCmd.Flags().BoolVar(&checksOpts.version, "version", false, "check Docker version")
	checkCmd.Flags().BoolVar(&checksOpts.scheduling, "scheduling", false, "check Docker service scheduling")
	checkCmd.Flags().BoolVarP(&checksOpts.all, "all", "a", false, "all tests")

	rootCmd.AddCommand(checkCmd)

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "set up amp services in swarm environment",
		Run:   install,
	}
	rootCmd.AddCommand(installCmd)

	_ = rootCmd.Execute()
}
