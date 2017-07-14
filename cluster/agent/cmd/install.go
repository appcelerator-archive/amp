package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/appcelerator/amp/cluster/agent/pkg/docker"
	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/pkg/term"
	"github.com/spf13/cobra"
	"github.com/subfuzion/stack/stack"
)

func NewInstallCommand() *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Set up amp services in swarm environment",
		RunE:  install,
	}
	return installCmd
}

func install(cmd *cobra.Command, args []string) error {
	stdin, stdout, stderr := term.StdStreams()
	dockerCli := docker.NewDockerCli(stdin, stdout, stderr)

	files, err := getStackFiles("./stacks")
	if err != nil {
		return err
	}
	for _, f := range files {
		log.Println(f)
		err := deploy(dockerCli, f)
		if err != nil {
			return err
		}
	}
	return nil
}

// returns sorted list of yaml file pathnames
func getStackFiles(path string) ([]string, error) {
	if path == "" {
		path = "./stacks"
	}

	// a bit more work but we can't just use filepath.Glob
	// since we need to match both *.yml and *.yaml
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	stackfiles := []string{}
	for _, f := range files {
		name := f.Name()
		// not compiling regex since only expecting less than a dozen stackfiles
		matched, err := regexp.MatchString("\\.ya?ml$", name)
		if err != nil {
			log.Println(err)
		} else if matched {
			stackfiles = append(stackfiles, filepath.Join(path, name))
		}
	}
	return stackfiles, nil
}

func deploy(d *command.DockerCli, stackfile string) error {
	// use the stackfile basename as the default stack namespace
	namespace := filepath.Base(stackfile)
	namespace = strings.TrimSuffix(namespace, filepath.Ext(namespace))

	opts := stack.DeployOptions{
		Namespace: namespace,
		Composefile: stackfile,
		ResolveImage: stack.ResolveImageNever,
		SendRegistryAuth: false,
		Prune: false,
	}
	err := stack.Deploy(d, opts)
	return err
}
