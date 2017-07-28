package main

import (
	"log"

	adm "github.com/appcelerator/amp/cluster/agent/admin"
	"github.com/spf13/cobra"
)

type CheckOptions struct {
	version    bool
	labels     bool
	scheduling bool
	all        bool
}

var checksOpts = &CheckOptions{}

func NewChecksCommand() *cobra.Command {
	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Run validation tests on the cluster",
		RunE:   checks,
	}
	checkCmd.Flags().BoolVar(&checksOpts.version, "version", false, "check Docker version")
	checkCmd.Flags().BoolVar(&checksOpts.labels, "labels", false, "check all required labels are defined on the swarm")
	checkCmd.Flags().BoolVar(&checksOpts.scheduling, "scheduling", false, "check Docker service scheduling")

	return checkCmd
}

func checks(cmd *cobra.Command, args []string) error {
	// if zero tests have been explicitely asked, run them all
	if !checksOpts.version && !checksOpts.labels && !checksOpts.scheduling {
		checksOpts.all = true
	}
	if checksOpts.version || checksOpts.all {
		if err := adm.VerifyDockerVersion(); err != nil {
			log.Println("Version test: FAIL")
			return err
		} else {
			log.Println("Version test: PASS")
		}
	}
	if checksOpts.labels || checksOpts.all {
		if err := adm.VerifyLabels(); err != nil {
			log.Println("Labels test: FAIL")
			return err
		} else {
			log.Println("Labels test: PASS")
		}
	}
	if checksOpts.scheduling || checksOpts.all {
		if err := adm.VerifyServiceScheduling(); err != nil {
			log.Println("Labels test: FAIL")
			return err
		} else {
			log.Println("Labels test: PASS")
		}
	}

	return nil
}
