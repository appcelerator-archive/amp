package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "envoy",
		Short: "run commands in target cluster",
		// If needed
		// PersistentPreRun: initAdmin,
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

	_ = rootCmd.Execute()
}
