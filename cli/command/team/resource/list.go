package resource

import (
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

type listTeamResOptions struct {
	org   string
	team  string
	quiet bool
}

// NewListTeamResCommand returns a new instance of the list team resource command.
func NewListTeamResCommand(c cli.Interface) *cobra.Command {
	opts := listTeamResOptions{}
	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Short:   "List resources",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTeamRes(c, cmd, opts)
		},
	}
	flags := cmd.Flags()
	//flags.StringVar(&opts.org, "org", "", "Organization name")
	flags.StringVar(&opts.team, "team", "", "Team name")
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only display team resource names")
	return cmd
}

func listTeamRes(c cli.Interface, cmd *cobra.Command, opts listTeamResOptions) error {
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
	client := account.NewAccountClient(conn)
	request := &account.GetTeamRequest{
		OrganizationName: opts.org,
		TeamName:         opts.team,
	}
	reply, err := client.GetTeam(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	//if err := cli.SaveOrg(opts.org, c.Server()); err != nil {
	//	return err
	//}
	if err := cli.SaveTeam(opts.team, c.Server()); err != nil {
		return err
	}
	if opts.quiet {
		for _, resource := range reply.Team.Resources {
			c.Console().Println(resource.Id)
		}
		return nil
	}
	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "RESOURCE ID\tPERMISSION LEVEL")
	for _, resource := range reply.Team.Resources {
		fmt.Fprintf(w, "%s\t%s\n", resource.Id, resource.PermissionLevel)
	}
	w.Flush()
	return nil
}
