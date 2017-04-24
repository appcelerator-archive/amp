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

type getTeamOptions struct {
	org  string
	team string
}

// NewTeamGetCommand returns a new instance of the get team command.
func NewTeamGetCommand(c cli.Interface) *cobra.Command {
	opts := getTeamOptions{}
	cmd := &cobra.Command{
		Use:     "get [OPTIONS]",
		Short:   "Get team",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getTeam(c, cmd, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.org, "org", "", "Organization name")
	flags.StringVar(&opts.team, "team", "", "Team name")
	return cmd
}

func getTeam(c cli.Interface, cmd *cobra.Command, opts getTeamOptions) error {
	if !cmd.Flag("org").Changed {
		opts.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		opts.team = c.Console().GetInput("team name")
	}

	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.GetTeamRequest{
		OrganizationName: opts.org,
		TeamName:         opts.team,
	}
	reply, err := client.GetTeam(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Printf("Team: %s\n", reply.Team.Name)
	c.Console().Printf("Organization: %s\n", opts.org)
	c.Console().Printf("Created On: %s\n", time.ConvertTime(reply.Team.CreateDt))
	return nil
}
