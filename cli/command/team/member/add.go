package member

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type addTeamMemOptions struct {
	org    string
	team   string
	member string
}

// NewAddTeamMemCommand returns a new instance of the add team member command.
func NewAddTeamMemCommand(c cli.Interface) *cobra.Command {
	opts := addTeamMemOptions{}
	cmd := &cobra.Command{
		Use:     "add [OPTIONS]",
		Short:   "Add member to team",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addTeamMem(c, cmd, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.org, "org", "", "Organization name")
	flags.StringVar(&opts.team, "team", "", "Team name")
	flags.StringVar(&opts.member, "member", "", "Member name")
	return cmd
}

func addTeamMem(c cli.Interface, cmd *cobra.Command, opts addTeamMemOptions) error {
	if !cmd.Flag("org").Changed {
		opts.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		opts.team = c.Console().GetInput("team name")
	}
	if !cmd.Flag("member").Changed {
		opts.member = c.Console().GetInput("member name")
	}

	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.AddUserToTeamRequest{
		OrganizationName: opts.org,
		TeamName:         opts.team,
		UserName:         opts.member,
	}
	if _, err := client.AddUserToTeam(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Member has been added to team.")
	return nil
}
