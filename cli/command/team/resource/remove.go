package resource

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type remTeamResOpts struct {
	org      string
	team     string
	resource string
}

var (
	remTeamResOptions = &remTeamResOpts{}
)

// NewRemoveTeamResCommand returns a new instance of the remove team resource command.
func NewRemoveTeamResCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm [OPTIONS]",
		Short:   "Remove resource from team",
		Aliases: []string{"del"},
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return remTeamRes(c, cmd)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&remTeamResOptions.org, "org", "", "Organization name")
	flags.StringVar(&remTeamResOptions.team, "team", "", "Team name")
	flags.StringVar(&remTeamResOptions.resource, "res", "", "Resource id")
	return cmd
}

func remTeamRes(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		remTeamResOptions.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		remTeamResOptions.team = c.Console().GetInput("team name")
	}
	if !cmd.Flag("res").Changed {
		remTeamResOptions.resource = c.Console().GetInput("resource id")
	}

	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.RemoveResourceFromTeamRequest{
		OrganizationName: remTeamResOptions.org,
		TeamName:         remTeamResOptions.team,
		ResourceId:       remTeamResOptions.resource,
	}
	if _, err := client.RemoveResourceFromTeam(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Resource has been removed from team.")
	return nil
}
