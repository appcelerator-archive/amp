package member

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type remTeamMemOpts struct {
	org    string
	team   string
	member string
}

var (
	remTeamMemOptions = &remTeamMemOpts{}
)

// NewRemoveTeamMemCommand returns a new instance of the remove team member command.
func NewRemoveTeamMemCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm",
		Short:   "Remove member from team",
		Aliases: []string{"del"},
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return remTeamMem(c, cmd)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&remTeamMemOptions.org, "org", "", "Organization name")
	flags.StringVar(&remTeamMemOptions.team, "team", "", "Team name")
	flags.StringVar(&remTeamMemOptions.team, "member", "", "Member name")
	return cmd
}

func remTeamMem(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		remTeamMemOptions.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		remTeamMemOptions.team = c.Console().GetInput("team name")
	}
	if !cmd.Flag("member").Changed {
		remTeamMemOptions.member = c.Console().GetInput("member name")
	}

	conn, err := c.ClientConn()
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.RemoveUserFromTeamRequest{
		OrganizationName: remTeamMemOptions.org,
		TeamName:         remTeamMemOptions.team,
		UserName:         remTeamMemOptions.member,
	}
	if _, err = client.RemoveUserFromTeam(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Member has been removed from team.")
	return nil
}
