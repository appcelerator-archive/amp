package team

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/pkg/time"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type getTeamOpts struct {
	org  string
	team string
}

var (
	getTeamOptions = &getTeamOpts{}
)

// NewTeamGetCommand returns a new instance of the get team command.
func NewTeamGetCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get [OPTIONS]",
		Short:   "Get team",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getTeam(c, cmd)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&getTeamOptions.org, "org", "", "Organization name")
	flags.StringVar(&getTeamOptions.team, "team", "", "Team name")
	return cmd
}

func getTeam(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		getTeamOptions.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		getTeamOptions.org = c.Console().GetInput("team name")
	}

	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.GetTeamRequest{
		OrganizationName: getTeamOptions.org,
		TeamName:         getTeamOptions.team,
	}
	reply, err := client.GetTeam(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Printf("Team: %s\n", reply.Team.Name)
	c.Console().Printf("Organization: %s\n", getTeamOptions.org)
	c.Console().Printf("Created On: %s\n", time.ConvertTime(reply.Team.CreateDt))
	return nil
}
