package resource

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type addTeamResOpts struct {
	org      string
	team     string
	resource string
}

var (
	addTeamResOptions = &addTeamResOpts{}
)

// NewAddTeamResCommand returns a new instance of the add team resource command.
func NewAddTeamResCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [OPTIONS]",
		Short:   "Add resource to team",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addTeamRes(c, cmd)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&addTeamResOptions.org, "org", "", "Organization name")
	flags.StringVar(&addTeamResOptions.team, "team", "", "Team name")
	flags.StringVar(&addTeamResOptions.resource, "res", "", "Resource id")
	return cmd
}

func addTeamRes(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		addTeamResOptions.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		addTeamResOptions.team = c.Console().GetInput("team name")
	}
	if !cmd.Flag("res").Changed {
		addTeamResOptions.resource = c.Console().GetInput("resource id")
	}

	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.AddResourceToTeamRequest{
		OrganizationName: addTeamResOptions.org,
		TeamName:         addTeamResOptions.team,
		ResourceId:       addTeamResOptions.resource,
	}
	if _, err := client.AddResourceToTeam(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Resource has been added to team.")
	return nil
}
