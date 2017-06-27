package member

import (
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

type listMemOrgOptions struct {
	name  string
	quiet bool
}

// NewOrgListMemCommand returns a new instance of the list organization member command.
func NewOrgListMemCommand(c cli.Interface) *cobra.Command {
	opts := listMemOrgOptions{}
	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Short:   "List members",
		Aliases: []string{"list"},
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listOrgMem(c, cmd, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.name, "org", "", "Organization name")
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only display member names")
	return cmd
}

func listOrgMem(c cli.Interface, cmd *cobra.Command, opts listMemOrgOptions) error {
	org, err := cli.ReadOrg(c.Server())
	if !cmd.Flag("org").Changed {
		switch {
		case err == nil:
			opts.name = org
			c.Console().Println("organization name:", opts.name)
		default:
			opts.name = c.Console().GetInput("organization name")
		}
	}
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.GetOrganizationRequest{
		Name: opts.name,
	}
	reply, err := client.GetOrganization(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	if err := cli.SaveOrg(opts.name, c.Server()); err != nil {
		return err
	}
	if opts.quiet {
		for _, member := range reply.Organization.Members {
			c.Console().Println(member.Name)
		}
		return nil
	}
	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "MEMBER\tROLE")
	for _, user := range reply.Organization.Members {
		fmt.Fprintf(w, "%s\t%s\n", user.Name, user.Role)
	}
	w.Flush()
	return nil
}
