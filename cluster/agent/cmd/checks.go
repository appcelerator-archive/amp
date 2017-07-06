package main

import (
	"log"

	adm "github.com/appcelerator/amp/cluster/agent/admin"
	"github.com/spf13/cobra"
)

type CheckOptions struct {
	version    bool
	scheduling bool
	all        bool
}

var checksOpts = &CheckOptions{}

func NewChecksCommand() *cobra.Command {
	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Run validation tests on the cluster",
		Run:   checks,
	}
	checkCmd.Flags().BoolVar(&checksOpts.version, "version", false, "check Docker version")
	checkCmd.Flags().BoolVar(&checksOpts.scheduling, "scheduling", false, "check Docker service scheduling")
	checkCmd.Flags().BoolVarP(&checksOpts.all, "all", "a", false, "all tests")

	return checkCmd
}

func checks(cmd *cobra.Command, args []string) {
	if checksOpts.version || checksOpts.all {
		out, err := adm.VerifyDockerVersion()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(out)
	}
	if checksOpts.scheduling || checksOpts.all {
		out, err := adm.VerifyServiceScheduling()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(out)
	}
}
