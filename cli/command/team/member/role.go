package member

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type changeMemTeamOpts struct {
	org    string
	team   string
	member string
	role   string
}

var (
	changeMemTeamOptions = &changeMemTeamOpts{}
)

// NewTeamChangeMemRoleCommand returns a new instance of the team member role change command.
func NewTeamChangeMemRoleCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "role",
		Short:   "Change member role",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return changeOrgMemRole(c, cmd)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&changeMemTeamOptions.org, "org", "", "Organization name")
	flags.StringVar(&changeMemTeamOptions.org, "team", "", "Team name")
	flags.StringVar(&changeMemTeamOptions.member, "member", "", "Member name")
	flags.StringVar(&changeMemTeamOptions.role, "role", "", "Member role")
	return cmd
}

func changeOrgMemRole(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		changeMemTeamOptions.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		changeMemTeamOptions.team = c.Console().GetInput("team name")
	}
	if !cmd.Flag("member").Changed {
		changeMemTeamOptions.member = c.Console().GetInput("member name")
	}
	if !cmd.Flag("role").Changed {
		changeMemTeamOptions.role = c.Console().GetInput("organization role")
	}
	teamRole := accounts.TeamRole_TEAM_MEMBER
	switch changeMemTeamOptions.role {
	case "owner":
		teamRole = accounts.TeamRole_TEAM_OWNER
	case "member":
		teamRole = accounts.TeamRole_TEAM_MEMBER
	default:
		return fmt.Errorf("invalid team role: %s. Please specify 'owner' or 'member' as role value.", changeMemTeamOptions.role)
	}
	conn, err := c.ClientConn()
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.ChangeTeamMemberRoleRequest{
		OrganizationName: changeMemTeamOptions.org,
		TeamName:         changeMemTeamOptions.team,
		UserName:         changeMemTeamOptions.member,
		Role:             teamRole,
	}
	if _, err = client.ChangeTeamMemberRole(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Member role has been changed.")
	return nil
}
