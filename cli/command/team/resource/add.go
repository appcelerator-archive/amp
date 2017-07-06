package resource

import (
	"errors"
	"strings"

	"github.com/appcelerator/amp/api/rpc/resource"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

type addTeamResOptions struct {
	org  string
	team string
}

// NewAddTeamResCommand returns a new instance of the add team resource command.
func NewAddTeamResCommand(c cli.Interface) *cobra.Command {
	opts := addTeamResOptions{}
	cmd := &cobra.Command{
		Use:     "add [OPTIONS] RESOURCE(S)",
		Short:   "Add one or more resources",
		PreRunE: cli.AtLeastArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addTeamRes(c, cmd, args, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.org, "org", "", "Organization name")
	flags.StringVar(&opts.team, "team", "", "Team name")
	return cmd
}

func addTeamRes(c cli.Interface, cmd *cobra.Command, args []string, opts addTeamResOptions) error {
	var errs []string
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

	conn := c.ClientConn()
	client := resource.NewResourceClient(conn)
	for _, res := range args {
		request := &resource.AddToTeamRequest{
			ResourceId:       res,
			OrganizationName: opts.org,
			TeamName:         opts.team,
		}
		if _, err := client.AddToTeam(context.Background(), request); err != nil {
			if s, ok := status.FromError(err); ok {
				errs = append(errs, s.Message())
				continue
			}
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	if err := cli.SaveOrg(opts.org, c.Server()); err != nil {
		return err
	}
	if err := cli.SaveTeam(opts.team, c.Server()); err != nil {
		return err
	}
	c.Console().Println("Resource(s) have been added to team.")
	return nil
}
