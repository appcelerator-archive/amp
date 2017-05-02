package resource

import (
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
	flags.StringVar(&opts.org, "org", "", "Organization name")
	flags.StringVar(&opts.team, "team", "", "Team name")
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only display team resource names")
	return cmd
}

func listTeamRes(c cli.Interface, cmd *cobra.Command, opts listTeamResOptions) error {
	if !cmd.Flag("org").Changed {
		opts.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		opts.team = c.Console().GetInput("team name")
	}
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.GetTeamRequest{
		OrganizationName: opts.org,
		TeamName:         opts.team,
	}
	reply, err := client.GetTeam(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
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
