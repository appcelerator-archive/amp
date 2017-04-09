package team

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type listTeamOpts struct {
	org   string
	quiet bool
}

var (
	listTeamOptions = &listTeamOpts{}
)

// NewTeamListCommand returns a new instance of the team list command.
func NewTeamListCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "List team",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTeam(c, cmd)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&listTeamOptions.org, "org", "", "Organization name")
	flags.BoolVarP(&listTeamOptions.quiet, "quiet", "q", false, "Only display team names")
	return cmd
}

func listTeam(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		listTeamOptions.org = c.Console().GetInput("organization name")
	}

	conn, err := c.ClientConn()
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.ListTeamsRequest{
		OrganizationName: listTeamOptions.org,
	}
	reply, err := client.ListTeams(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if listTeamOptions.quiet {
		for _, team := range reply.Teams {
			c.Console().Println(team.Name)
		}
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', 0)
	fmt.Fprintln(w, "TEAM\tCREATED\t")
	for _, team := range reply.Teams {
		fmt.Fprintf(w, "%s\t%s\t\n", team.Name, cli.ConvertTime(c, team.CreateDt))
	}
	w.Flush()
	return nil
}
