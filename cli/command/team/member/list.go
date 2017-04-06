package member

import (
	"fmt"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"text/tabwriter"
)

type listTeamMemOpts struct {
	org  string
	team string
}

var (
	listTeamMemOptions = &listTeamMemOpts{}
)

// NewListTeamMemCommand returns a new instance of the list team member command.
func NewListTeamMemCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "List members",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTeamMem(c, cmd)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&listTeamMemOptions.org, "org", "", "Organization name")
	flags.StringVar(&listTeamMemOptions.team, "team", "", "Team name")
	flags.BoolP("quiet", "q", false, "Only display team member names")
	return cmd
}

func listTeamMem(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		listTeamMemOptions.org = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("team").Changed {
		listTeamMemOptions.team = c.Console().GetInput("team name")
	}

	conn, err := c.ClientConn()
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.GetTeamRequest{
		OrganizationName: listTeamMemOptions.org,
		TeamName:         listTeamMemOptions.team,
	}
	reply, err := client.GetTeam(context.Background(), request)
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	if quiet, err := strconv.ParseBool(cmd.Flag("quiet").Value.String()); err != nil {
		c.Console().Fatalf("unable to convert quiet parameter : %v", err.Error())
	} else if quiet {
		for _, member := range reply.Team.Members {
			c.Console().Println(member.Name)
		}
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', 0)
	fmt.Fprintln(w, "TEAM\tROLE\t")
	for _, team := range reply.Team.Members {
		fmt.Fprintf(w, "%s\t%s\n", team.Name, team.Role)
	}
	w.Flush()
	return nil
}
