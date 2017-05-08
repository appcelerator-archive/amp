package team

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type createTeamOptions struct {
	org  string
	team string
}

// NewTeamCreateCommand returns a new instance of the team create command.
func NewTeamCreateCommand(c cli.Interface) *cobra.Command {
	opts := createTeamOptions{}
	cmd := &cobra.Command{
		Use:     "create [OPTIONS]",
		Short:   "Create team",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createTeam(c, cmd, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.org, "org", "", "Organization name")
	flags.StringVar(&opts.team, "team", "", "Team name")
	return cmd
}

func createTeam(c cli.Interface, cmd *cobra.Command, opts createTeamOptions) error {
	if !cmd.Flag("org").Changed {
		opts.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		opts.team = c.Console().GetInput("team name")
	}

	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.CreateTeamRequest{
		OrganizationName: opts.org,
		TeamName:         opts.team,
	}
	if _, err := client.CreateTeam(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Team has been created in the organization.")
	return nil
}
