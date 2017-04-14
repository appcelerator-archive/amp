package resource

import (
	"errors"
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type changeTeamResPermLevelOpts struct {
	org             string
	team            string
	resource        string
	permissionLevel string
}

var (
	changeTeamResPermLevelOptions = &changeTeamResPermLevelOpts{}
)

// NewChangeTeamResPermissionLevelCommand returns a new instance of the team resource perm command.
func NewChangeTeamResPermissionLevelCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "perm [OPTIONS] RESOURCEID PERMISSION",
		Short:   "Change permission level over a resource",
		PreRunE: cli.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return errors.New("resource id cannot be empty")
			}
			if args[1] == "" {
				return errors.New("permission level cannot be empty")
			}
			changeTeamResPermLevelOptions.resource = args[0]
			changeTeamResPermLevelOptions.permissionLevel = args[1]
			return changeOrgMemRole(c, cmd)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&changeTeamResPermLevelOptions.org, "org", "", "Organization name")
	flags.StringVar(&changeTeamResPermLevelOptions.team, "team", "", "Team name")
	return cmd
}

func changeOrgMemRole(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		changeTeamResPermLevelOptions.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		changeTeamResPermLevelOptions.team = c.Console().GetInput("team name")
	}
	permissionLevel := accounts.TeamPermissionLevel_TEAM_READ
	switch changeTeamResPermLevelOptions.permissionLevel {
	case "read":
		permissionLevel = accounts.TeamPermissionLevel_TEAM_READ
	case "write":
		permissionLevel = accounts.TeamPermissionLevel_TEAM_WRITE
	case "admin":
		permissionLevel = accounts.TeamPermissionLevel_TEAM_ADMIN
	default:
		return fmt.Errorf("invalid permission level: %s. Please specify 'read', 'write' or 'admin' as permission value.", changeTeamResPermLevelOptions.permissionLevel)
	}
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.ChangeTeamResourcePermissionLevelRequest{
		OrganizationName: changeTeamResPermLevelOptions.org,
		TeamName:         changeTeamResPermLevelOptions.team,
		ResourceId:       changeTeamResPermLevelOptions.resource,
		PermissionLevel:  permissionLevel,
	}
	if _, err := client.ChangeTeamResourcePermissionLevel(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Permission level has been changed.")
	return nil
}
