package cmd

import (
	"log"

	"fmt"

	"docker.io/go-docker"
	"docker.io/go-docker/api"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/versions"
	ampdocker "github.com/appcelerator/amp/pkg/docker"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
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
		RunE:  Checks,
	}
	checkCmd.Flags().BoolVar(&checksOpts.version, "version", false, "check Docker version")
	checkCmd.Flags().BoolVar(&checksOpts.labels, "labels", false, "check all required labels are defined on the swarm")
	checkCmd.Flags().BoolVar(&checksOpts.scheduling, "scheduling", false, "check Docker service scheduling")

	return checkCmd
}

func Checks(cmd *cobra.Command, args []string) error {
	// if zero tests have been explicitly asked, run them all
	if !checksOpts.version && !checksOpts.labels && !checksOpts.scheduling {
		checksOpts.all = true
	}
	if checksOpts.version || checksOpts.all {
		if err := VerifyDockerVersion(); err != nil {
			log.Println("Version test: FAIL")
			return err
		} else {
			log.Println("Version test: PASS")
		}
	}
	if checksOpts.labels || checksOpts.all {
		if err := VerifyLabels(); err != nil {
			log.Println("Labels test: FAIL")
			return err
		} else {
			log.Println("Labels test: PASS")
		}
	}
	//if checksOpts.scheduling || checksOpts.all {
	//	if err := adm.VerifyServiceScheduling(); err != nil {
	//		log.Println("Service scheduling test: FAIL")
	//		return err
	//	} else {
	//		log.Println("Service scheduling test: PASS")
	//	}
	//}

	return nil
}

func VerifyDockerVersion() error {
	c, err := docker.NewClient(ampdocker.DefaultURL, ampdocker.DefaultVersion, nil, nil)
	if err != nil {
		return err
	}
	version, err := c.ServerVersion(context.Background())
	apiVersion := version.APIVersion
	if versions.LessThan(apiVersion, ampdocker.MinVersion) {
		log.Printf("Docker engine version %s\n", version.Version)
		log.Printf("API version - minimum expected: %.s, observed: %.s", ampdocker.MinVersion, apiVersion)
		return errors.New("Docker engine doesn't meet the requirements (API Version)")
	}
	return nil
}

func VerifyLabels() error {
	labels := map[string]bool{}
	expectedLabels := []string{"amp.type.api=true", "amp.type.route=true", "amp.type.core=true", "amp.type.metrics=true",
		"amp.type.search=true", "amp.type.mq=true", "amp.type.kv=true", "amp.type.user=true"}
	missingLabel := false
	c, err := docker.NewClient(ampdocker.DefaultURL, api.DefaultVersion, nil, nil)
	if err != nil {
		return err
	}
	nodes, err := c.NodeList(context.Background(), types.NodeListOptions{})
	if err != nil {
		return err
	}
	// get the full list of labels
	for _, node := range nodes {
		nodeLabels := node.Spec.Annotations.Labels
		for k, v := range nodeLabels {
			labels[fmt.Sprintf("%s=%s", k, v)] = true
		}
	}
	// check that all expected labels are at least on one node
	for _, label := range expectedLabels {
		if !labels[label] {
			log.Printf("label %s is missing\n", label)
			missingLabel = true
		}
	}
	if missingLabel {
		return errors.New("At least one missing label")
	}
	return nil

}
