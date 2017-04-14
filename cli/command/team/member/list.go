package member

import (
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type listTeamMemOpts struct {
	org   string
	team  string
	quiet bool
}

var (
	listTeamMemOptions = &listTeamMemOpts{}
)

// NewListTeamMemCommand returns a new instance of the list team member command.
func NewListTeamMemCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Short:   "List members",
		Aliases: []string{"list"},
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTeamMem(c, cmd)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&listTeamMemOptions.org, "org", "", "Organization name")
	flags.StringVar(&listTeamMemOptions.team, "team", "", "Team name")
	flags.BoolVarP(&listTeamMemOptions.quiet, "quiet", "q", false, "Only display team member names")
	return cmd
}

func listTeamMem(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		listTeamMemOptions.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		listTeamMemOptions.team = c.Console().GetInput("team name")
	}

	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.GetTeamRequest{
		OrganizationName: listTeamMemOptions.org,
		TeamName:         listTeamMemOptions.team,
	}
	reply, err := client.GetTeam(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if listTeamMemOptions.quiet {
		for _, member := range reply.Team.Members {
			c.Console().Println(member.Name)
		}
		return nil
	}
	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "MEMBER")
	for _, member := range reply.Team.Members {
		fmt.Fprintf(w, "%s\n", member)
	}
	w.Flush()
	return nil
}
