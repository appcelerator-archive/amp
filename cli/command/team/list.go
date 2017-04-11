package team

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/pkg/time"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type listTeamOpts struct {
	org string
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
	flags.BoolP("quiet", "q", false, "Only display team names")
	return cmd
}

func listTeam(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		listTeamOptions.org = c.Console().GetInput("organization name")
	}

	conn, err := c.ClientConn()
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.ListTeamsRequest{
		OrganizationName: listTeamOptions.org,
	}
	reply, err := client.ListTeams(context.Background(), request)
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	if quiet, err := strconv.ParseBool(cmd.Flag("quiet").Value.String()); err != nil {
		c.Console().Fatalf("unable to convert quiet parameter : %v", grpc.ErrorDesc(err))
	} else if quiet {
		for _, team := range reply.Teams {
			c.Console().Println(team.Name)
		}
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', 0)
	fmt.Fprintln(w, "TEAM\tCREATED\t")
	for _, team := range reply.Teams {
		fmt.Fprintf(w, "%s\t%s\t\n", team.Name, time.ConvertTime(team.CreateDt))
	}
	w.Flush()
	return nil
}
