package resource

import (
	"errors"
	"strings"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type remTeamResOptions struct {
	org       string
	team      string
	resources []string
}

// NewRemoveTeamResCommand returns a new instance of the remove team resource command.
func NewRemoveTeamResCommand(c cli.Interface) *cobra.Command {
	opts := remTeamResOptions{}
	cmd := &cobra.Command{
		Use:     "rm [OPTIONS] RESOURCE(S)",
		Short:   "Remove resource from team",
		Aliases: []string{"del"},
		PreRunE: cli.AtLeastArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return remTeamRes(c, cmd, args, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.org, "org", "", "Organization name")
	flags.StringVar(&opts.team, "team", "", "Team name")
	return cmd
}

func remTeamRes(c cli.Interface, cmd *cobra.Command, args []string, opts remTeamResOptions) error {
	var errs []string
	if !cmd.Flag("org").Changed {
		opts.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		opts.team = c.Console().GetInput("team name")
	}
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	for _, res := range args {
		request := &account.RemoveResourceFromTeamRequest{
			OrganizationName: opts.org,
			TeamName:         opts.team,
			ResourceId:       res,
		}
		if _, err := client.RemoveResourceFromTeam(context.Background(), request); err != nil {
			errs = append(errs, grpc.ErrorDesc(err))
			continue
		}
		c.Console().Println(res)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
