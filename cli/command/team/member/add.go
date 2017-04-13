package member

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type addTeamMemOpts struct {
	org    string
	team   string
	member string
}

var (
	addTeamMemOptions = &addTeamMemOpts{}
)

// NewAddTeamMemCommand returns a new instance of the add team member command.
func NewAddTeamMemCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [OPTIONS]",
		Short:   "Add member to team",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addTeamMem(c, cmd)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&addTeamMemOptions.org, "org", "", "Organization name")
	flags.StringVar(&addTeamMemOptions.team, "team", "", "Team name")
	flags.StringVar(&addTeamMemOptions.team, "member", "", "Member name")
	return cmd
}

func addTeamMem(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		addTeamMemOptions.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		addTeamMemOptions.team = c.Console().GetInput("team name")
	}
	if !cmd.Flag("member").Changed {
		addTeamMemOptions.member = c.Console().GetInput("member name")
	}

	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.AddUserToTeamRequest{
		OrganizationName: addTeamMemOptions.org,
		TeamName:         addTeamMemOptions.team,
		UserName:         addTeamMemOptions.member,
	}
	if _, err := client.AddUserToTeam(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Member has been added to team.")
	return nil
}
