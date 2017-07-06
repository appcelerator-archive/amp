package team

import (
	"errors"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/pkg/time"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

type getTeamOptions struct {
	org string
}

// NewTeamGetCommand returns a new instance of the get team command.
func NewTeamGetCommand(c cli.Interface) *cobra.Command {
	opts := getTeamOptions{}
	cmd := &cobra.Command{
		Use:     "get [OPTIONS] TEAM",
		Short:   "Get team information",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getTeam(c, cmd, args, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.org, "org", "", "Organization name")
	return cmd
}

func getTeam(c cli.Interface, cmd *cobra.Command, args []string, opts getTeamOptions) error {
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

	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.GetTeamRequest{
		OrganizationName: opts.org,
		TeamName:         args[0],
	}
	reply, err := client.GetTeam(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	c.Console().Printf("Team: %s\n", reply.Team.Name)
	c.Console().Printf("Organization: %s\n", opts.org)
	c.Console().Printf("Created On: %s\n", time.ConvertTime(reply.Team.CreateDt))
	return nil
}
