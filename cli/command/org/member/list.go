package member

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
		Use:     "ls",
		Short:   "List members",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listOrgMem(c, cmd)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&listMemOrgOptions.name, "org", "", "Organization name")
	flags.BoolVarP(&listMemOrgOptions.quiet, "quiet", "q", false, "Only display member names")
	return cmd
}

func listOrgMem(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		listMemOrgOptions.name = c.Console().GetInput("organization name")
	}
	conn, err := c.ClientConn()
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.GetOrganizationRequest{
		Name: listMemOrgOptions.name,
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
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', 0)
	fmt.Fprintln(w, "USERNAME\tROLE\t")
	for _, user := range reply.Organization.Members {
		fmt.Fprintf(w, "%s\t%s\n", user.Name, user.Role)
	}
	w.Flush()
	return nil
}
