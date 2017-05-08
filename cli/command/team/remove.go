package team

import (
	"errors"
	"strings"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type removeTeamOptions struct {
	org   string
	teams []string
}

// NewTeamRemoveCommand returns a new instance of the team remove command.
func NewTeamRemoveCommand(c cli.Interface) *cobra.Command {
	opts := removeTeamOptions{}
	cmd := &cobra.Command{
		Use:     "rm [OPTIONS] TEAM(S)",
		Short:   "Remove team",
		Aliases: []string{"remove"},
		PreRunE: cli.AtLeastArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeTeam(c, cmd, args, opts)
		},
	}
	cmd.Flags().StringVar(&opts.org, "org", "", "Organization name")
	return cmd
}

func removeTeam(c cli.Interface, cmd *cobra.Command, args []string, opts removeTeamOptions) error {
	var errs []string
	if !cmd.Flag("org").Changed {
		opts.org = c.Console().GetInput("organization name")
	}
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	for _, team := range args {
		request := &account.DeleteTeamRequest{
			OrganizationName: opts.org,
			TeamName:         team,
		}
		if _, err := client.DeleteTeam(context.Background(), request); err != nil {
			errs = append(errs, grpc.ErrorDesc(err))
			continue
		}
		c.Console().Println(team)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
