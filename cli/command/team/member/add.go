package member

import (
	"errors"
	"strings"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

type addTeamMemOptions struct {
	org  string
	team string
}

// NewAddTeamMemCommand returns a new instance of the add team member command.
func NewAddTeamMemCommand(c cli.Interface) *cobra.Command {
	opts := addTeamMemOptions{}
	cmd := &cobra.Command{
		Use:     "add [OPTIONS] MEMBER(S)",
		Short:   "Add one or more members",
		PreRunE: cli.AtLeastArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addTeamMem(c, cmd, args, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.org, "org", "", "Organization name")
	flags.StringVar(&opts.team, "team", "", "Team name")
	return cmd
}

func addTeamMem(c cli.Interface, cmd *cobra.Command, args []string, opts addTeamMemOptions) error {
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
	client := account.NewAccountClient(conn)
	for _, member := range args {
		request := &account.AddUserToTeamRequest{
			OrganizationName: opts.org,
			TeamName:         opts.team,
			UserName:         member,
		}
		if _, err := client.AddUserToTeam(context.Background(), request); err != nil {
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
	c.Console().Println("Member(s) have been added to team.")
	return nil
}
