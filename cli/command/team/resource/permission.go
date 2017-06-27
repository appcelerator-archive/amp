package resource

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/resource"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type changeTeamResPermLevelOptions struct {
	org             string
	team            string
	resource        string
	permissionLevel string
}

// NewChangeTeamResPermissionLevelCommand returns a new instance of the team resource perm command.
func NewChangeTeamResPermissionLevelCommand(c cli.Interface) *cobra.Command {
	opts := changeTeamResPermLevelOptions{}
	cmd := &cobra.Command{
		Use:     "perm [OPTIONS] RESOURCEID PERMISSION",
		Short:   "Change permission level over a resource",
		PreRunE: cli.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.resource = args[0]
			opts.permissionLevel = args[1]
			return changeOrgMemRole(c, cmd, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.org, "org", "", "Organization name")
	flags.StringVar(&opts.team, "team", "", "Team name")
	return cmd
}

func changeOrgMemRole(c cli.Interface, cmd *cobra.Command, opts changeTeamResPermLevelOptions) error {
	org, err := cli.ReadOrg(c.Server())
	if !cmd.Flag("org").Changed {
		switch {
		case err == nil:
			opts.org = org
			c.Console().Println("organization name:", opts.org)
		default:
			opts.org = c.Console().GetInput("organization name")
		}
	}
	team, err := cli.ReadTeam(c.Server())
	if !cmd.Flag("team").Changed {
		switch {
		case err == nil:
			opts.team = team
			c.Console().Println("team name:", opts.team)
		default:
			opts.team = c.Console().GetInput("team name")
		}
	}
	permissionLevel := accounts.TeamPermissionLevel_TEAM_READ
	switch opts.permissionLevel {
	case "read":
		permissionLevel = accounts.TeamPermissionLevel_TEAM_READ
	case "write":
		permissionLevel = accounts.TeamPermissionLevel_TEAM_WRITE
	case "admin":
		permissionLevel = accounts.TeamPermissionLevel_TEAM_ADMIN
	default:
		return fmt.Errorf("invalid permission level: %s. Please specify 'read', 'write' or 'admin' as permission value.", opts.permissionLevel)
	}
	conn := c.ClientConn()
	client := resource.NewResourceClient(conn)
	request := &resource.ChangePermissionLevelRequest{
		OrganizationName: opts.org,
		TeamName:         opts.team,
		ResourceId:       opts.resource,
		PermissionLevel:  permissionLevel,
	}
	if _, err := client.ChangePermissionLevel(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if err := cli.SaveOrg(opts.org, c.Server()); err != nil {
		return err
	}
	if err := cli.SaveTeam(opts.team, c.Server()); err != nil {
		return err
	}
	c.Console().Println("Permission level has been changed.")
	return nil
}
