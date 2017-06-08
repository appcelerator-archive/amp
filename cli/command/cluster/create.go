package cluster

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/spf13/cobra"
)

// NewCreateCommand returns a new instance of the create command for bootstrapping a cluster.
func NewCreateCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [OPTIONS]",
		Short:   "Create an amp cluster",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(c, cmd)
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&opts.log, "log", "L", 4, "Logging level. 0 is least verbose, Max is 5")
	flags.StringVarP(&opts.name, "name", "i", "", "Cluster label used as the cluster id")
	flags.StringVarP(&opts.tag, "tag", "t", c.Version(), "Tag for cluster images, use 'local' for development")
	flags.StringVarP(&opts.provider, "provider", "p", "local", "Cluster provider, options are 'local' and 'docker'")
	flags.IntVarP(&opts.managers, "managers", "m", 3, "Initial number of manager nodes")
	flags.IntVarP(&opts.workers, "workers", "w", 2, "Initial number of worker nodes")
	flags.StringVarP(&opts.registration, "registration", "r", configuration.RegistrationNone, "Registration policy, options are 'none' or 'email'")
	flags.BoolVarP(&opts.notifications, "notifications", "n", false, "Enable/disable server notifications")
	flags.StringVarP(&opts.region, "region", "g", "", "Region for deployment on selected cloud provider")
	flags.StringVarP(&opts.domain, "domain", "d", "", "ssh key for deployment on selected cloud provider")
	return cmd
}

// Map cli cluster flags to target bootstrap cluster command flags and update the cluster
func create(c cli.Interface, cmd *cobra.Command) error {
	// This is a map from cli cluster flag name to bootstrap script flag name
	m := map[string]string{
		"log":           "-L",
		"name":          "-i",
		"tag":           "-T",
		"provider":      "-p",
		"managers":      "-m",
		"workers":       "-w",
		"registration":  "-r",
		"notifications": "-n",
		"region":        "-g",
		"domain":        "-D",
	}

	switch opts.provider {
	case "local", "docker":
		// Check that server is local
		local := strings.HasPrefix(c.Server(), "127.0.0.1") || strings.HasPrefix(c.Server(), "localhost")
		if !local {
			return fmt.Errorf("can only deploy a %s cluster on '--server=localhost'", opts.provider)
		}
		// Ignore certain flags for different deployments
		// cluster size and cloud vars
		// TODO: docker infrakit plugin should not use static managers and workers
		delete(m, "managers")
		delete(m, "workers")
		delete(m, "region")
		delete(m, "domain")
	case "aws":
		// Specific case handling for cloud providers that will be supported in future
		return fmt.Errorf("%s cluster deployment is not supported in this release: %v", opts.provider, c.Version())
	default:
		// default should be used for any cloud deployment as it will take all flags
		// for now, this is not supported
		return fmt.Errorf("%s cluster deployment is not supported", opts.provider)
	}

	args := []string{"bin/deploy"}
	args = reflag(cmd, m, args)
	env := map[string]string{"TAG": opts.tag, "REGISTRATION": opts.registration, "NOTIFICATIONS": strconv.FormatBool(opts.notifications)}
	err := queryCluster(c, args, env)
	if err != nil {
		return err
	}
	c.Console().Println("")
	c.Console().Successln("See `amp user signup --help` in order to create your first user.")
	c.Console().Successln("Read the documentation here https://github.com/appcelerator/amp/tree/master/docs.")
	// TODO: point to cloud.appcelerator.io for cloud deployments
	c.Console().Successln("The web portal is accessible at https://local.appcelerator.io/")
	return err
}
