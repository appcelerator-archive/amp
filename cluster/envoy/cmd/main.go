package main

import (
	"log"

	"github.com/appcelerator/amp/cluster/envoy"
	"github.com/spf13/cobra"
)

type CheckOptions struct {
	version    bool
	scheduling bool
	all        bool
}

var checksOpts = &CheckOptions{}

func checks(cmd *cobra.Command, args []string) {
	if checksOpts.version || checksOpts.all {
		out, err := envoy.VerifyDockerVersion()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(out)
	}
	if checksOpts.scheduling || checksOpts.all {
		out, err := envoy.VerifyServiceScheduling()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(out)
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "envoy",
		Short: "run commands in target cluster",
		// If needed
		// PersistentPreRun: initEnvoy,
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
