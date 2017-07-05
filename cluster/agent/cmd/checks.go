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
