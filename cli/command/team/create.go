package team

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type createTeamOpts struct {
	org  string
	team string
}

var (
	createTeamOptions = &createTeamOpts{}
)

// NewTeamCreateCommand returns a new instance of the team create command.
func NewTeamCreateCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [OPTIONS]",
		Short:   "Create team",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createTeam(c, cmd)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&createTeamOptions.org, "org", "", "Organization name")
	flags.StringVar(&createTeamOptions.team, "team", "", "Team name")
	return cmd
}

func createTeam(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		createTeamOptions.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		createTeamOptions.team = c.Console().GetInput("team name")
	}

	conn, err := c.ClientConn()
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.CreateTeamRequest{
		OrganizationName: createTeamOptions.org,
		TeamName:         createTeamOptions.team,
	}
	if _, err = client.CreateTeam(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Team has been created in the organization.")
	return nil
}
