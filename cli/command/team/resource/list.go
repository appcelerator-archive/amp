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

type listTeamResOpts struct {
	org   string
	team  string
	quiet bool
}

var (
	listTeamResOptions = &listTeamResOpts{}
)

// NewListTeamResCommand returns a new instance of the list team resource command.
func NewListTeamResCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Short:   "List resources",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTeamRes(c, cmd)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&listTeamResOptions.org, "org", "", "Organization name")
	flags.StringVar(&listTeamResOptions.team, "team", "", "Team name")
	flags.BoolVarP(&listTeamResOptions.quiet, "quiet", "q", false, "Only display team resource names")
	return cmd
}

func listTeamRes(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		listTeamResOptions.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		listTeamResOptions.team = c.Console().GetInput("team name")
	}

	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.GetTeamRequest{
		OrganizationName: listTeamResOptions.org,
		TeamName:         listTeamResOptions.team,
	}
	reply, err := client.GetTeam(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if listTeamResOptions.quiet {
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
