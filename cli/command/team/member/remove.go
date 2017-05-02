package member

import (
	"errors"
	"strings"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type remTeamMemOptions struct {
	org  string
	team string
}

// NewRemoveTeamMemCommand returns a new instance of the remove team member command.
func NewRemoveTeamMemCommand(c cli.Interface) *cobra.Command {
	opts := remTeamMemOptions{}
	cmd := &cobra.Command{
		Use:     "rm [OPTIONS] MEMBER(S)",
		Short:   "Remove member from team",
		Aliases: []string{"remove"},
		PreRunE: cli.AtLeastArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return remTeamMem(c, cmd, args, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.org, "org", "", "Organization name")
	flags.StringVar(&opts.team, "team", "", "Team name")
	return cmd
}

func remTeamMem(c cli.Interface, cmd *cobra.Command, args []string, opts remTeamMemOptions) error {
	var errs []string
	if !cmd.Flag("org").Changed {
		opts.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		opts.team = c.Console().GetInput("team name")
	}
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	for _, member := range args {
		request := &account.RemoveUserFromTeamRequest{
			OrganizationName: opts.org,
			TeamName:         opts.team,
			UserName:         member,
		}
		if _, err := client.RemoveUserFromTeam(context.Background(), request); err != nil {
			errs = append(errs, grpc.ErrorDesc(err))
			continue
		}
		c.Console().Println(member)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
