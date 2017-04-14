package member

import (
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type listMemOrgOpts struct {
	name  string
	quiet bool
}

var (
	listMemOrgOptions = &listMemOrgOpts{}
)

// NewOrgListMemCommand returns a new instance of the list organization member command.
func NewOrgListMemCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls [OPTIONS] ORGANIZATION",
		Short:   "List members",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return errors.New("organization name cannot be empty")
			}
			listMemOrgOptions.name = args[0]
			return listOrgMem(c, listMemOrgOptions)
		},
	}
	cmd.Flags().BoolVarP(&listMemOrgOptions.quiet, "quiet", "q", false, "Only display member names")
	return cmd
}

func listOrgMem(c cli.Interface, opt *listMemOrgOpts) error {
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.GetOrganizationRequest{
		Name: opt.name,
	}
	reply, err := client.GetOrganization(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if listMemOrgOptions.quiet {
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
