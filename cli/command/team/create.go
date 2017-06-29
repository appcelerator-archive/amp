package team

import (
	"errors"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

type createTeamOptions struct {
	org string
}

// NewTeamCreateCommand returns a new instance of the team create command.
func NewTeamCreateCommand(c cli.Interface) *cobra.Command {
	opts := createTeamOptions{}
	cmd := &cobra.Command{
		Use:     "create [OPTIONS] TEAM",
		Short:   "Create team",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createTeam(c, cmd, args, opts)
		},
	}
	cmd.Flags().StringVar(&opts.org, "org", "", "Organization name")
	return cmd
}

func createTeam(c cli.Interface, cmd *cobra.Command, args []string, opts createTeamOptions) error {
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

	team := args[0]
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.CreateTeamRequest{
		OrganizationName: opts.org,
		TeamName:         team,
	}
	if _, err := client.CreateTeam(context.Background(), request); err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	if err := cli.SaveOrg(opts.org, c.Server()); err != nil {
		return err
	}
	if err := cli.SaveTeam(team, c.Server()); err != nil {
		return err
	}
	c.Console().Println("Team has been created in the organization.")
	return nil
}
