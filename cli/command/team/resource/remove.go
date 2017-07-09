package resource

import (
	"errors"
	"strings"

	"github.com/appcelerator/amp/api/rpc/resource"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
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
		Short:   "Remove one or more resources",
		Aliases: []string{"del"},
		PreRunE: cli.AtLeastArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return remTeamRes(c, cmd, args, opts)
		},
	}
	flags := cmd.Flags()
	//flags.StringVar(&opts.org, "org", "", "Organization name")
	flags.StringVar(&opts.team, "team", "", "Team name")
	return cmd
}

func remTeamRes(c cli.Interface, cmd *cobra.Command, args []string, opts remTeamResOptions) error {
	var errs []string
	opts.org = accounts.DefaultOrganization
	//org, err := cli.ReadOrg(c.Server())
	//if !cmd.Flag("org").Changed {
	//	switch {
	//	case err == nil:
	//		opts.org = org
	//		c.Console().Println("organization name:", opts.org)
	//	default:
	//		opts.org = c.Console().GetInput("organization name")
	//	}
	//}
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
	conn := c.ClientConn()
	client := resource.NewResourceClient(conn)
	for _, res := range args {
		request := &resource.RemoveFromTeamRequest{
			ResourceId:       res,
			OrganizationName: opts.org,
			TeamName:         opts.team,
		}
		if _, err := client.RemoveFromTeam(context.Background(), request); err != nil {
			if s, ok := status.FromError(err); ok {
				errs = append(errs, s.Message())
				continue
			}
		}
		c.Console().Println(res)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	//if err := cli.SaveOrg(opts.org, c.Server()); err != nil {
	//	return err
	//}
	if err := cli.SaveTeam(opts.team, c.Server()); err != nil {
		return err
	}
	return nil
}
