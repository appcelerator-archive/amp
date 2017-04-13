package team

import (
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/pkg/time"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type listTeamOpts struct {
	name  string
	quiet bool
}

var (
	listTeamOptions = &listTeamOpts{}
)

// NewTeamListCommand returns a new instance of the team list command.
func NewTeamListCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls [OPTIONS] ORGANIZATION",
		Short:   "List team",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return errors.New("organization name cannot be empty")
			}
			listTeamOptions.name = args[0]
			return listTeam(c, listTeamOptions)
		},
	}
	cmd.Flags().BoolVarP(&listTeamOptions.quiet, "quiet", "q", false, "Only display team names")
	return cmd
}

func listTeam(c cli.Interface, opt *listTeamOpts) error {
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.ListTeamsRequest{
		OrganizationName: opt.name,
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
	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "TEAM\tCREATED ON")
	for _, team := range reply.Teams {
		fmt.Fprintf(w, "%s\t%s\t\n", team.Name, time.ConvertTime(team.CreateDt))
	}
	w.Flush()
	return nil
}
