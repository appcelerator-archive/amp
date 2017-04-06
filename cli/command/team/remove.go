package team

import (
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type removeTeamOpts struct {
	org  string
	team string
}

var (
	removeTeamOptions = &removeTeamOpts{}
)

// NewTeamRemoveCommand returns a new instance of the team remove command.
func NewTeamRemoveCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm",
		Short:   "Remove team",
		Aliases: []string{"del"},
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeTeam(c, cmd)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&removeTeamOptions.org, "org", "", "Organization name")
	flags.StringVar(&removeTeamOptions.team, "team", "", "Team name")
	return cmd
}

func removeTeam(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		removeTeamOptions.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		removeTeamOptions.team = c.Console().GetInput("team name")
	}

	conn, err := c.ClientConn()
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.DeleteTeamRequest{
		OrganizationName: removeTeamOptions.org,
		TeamName:         removeTeamOptions.team,
	}
	if _, err = client.DeleteTeam(context.Background(), request); err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	c.Console().Println("Team has been removed from the organization.")
	return nil
}
