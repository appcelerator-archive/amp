package resource

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/resource"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type addTeamResOptions struct {
	org      string
	team     string
	resource string
}

// NewAddTeamResCommand returns a new instance of the add team resource command.
func NewAddTeamResCommand(c cli.Interface) *cobra.Command {
	opts := addTeamResOptions{}
	cmd := &cobra.Command{
		Use:     "add [OPTIONS]",
		Short:   "Add resource to team",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addTeamRes(c, cmd, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.org, "org", "", "Organization name")
	flags.StringVar(&opts.team, "team", "", "Team name")
	flags.StringVar(&opts.resource, "res", "", "Resource id")
	return cmd
}

func addTeamRes(c cli.Interface, cmd *cobra.Command, opts addTeamResOptions) error {
	if !cmd.Flag("org").Changed {
		opts.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		opts.team = c.Console().GetInput("team name")
	}
	if !cmd.Flag("res").Changed {
		opts.resource = c.Console().GetInput("resource id")
	}

	conn := c.ClientConn()
	client := resource.NewResourceClient(conn)
	request := &resource.AddToTeamRequest{
		ResourceId:       opts.resource,
		OrganizationName: opts.org,
		TeamName:         opts.team,
	}
	if _, err := client.AddToTeam(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Resource has been added to team.")
	return nil
}
