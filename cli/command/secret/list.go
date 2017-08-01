package secret

import (
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/secret"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

type ListOpts struct {
	Filters []string
	Format  string
	Quiet   bool
}

var listOpts = &ListOpts{
	Filters: []string{},
}

// NewListCommand returns a new instance of the list command for listing secrets
func NewListCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [OPTIONS]",
		Short:   "List secrets",
		PreRunE: cli.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(c, cmd)
		},
	}

	flags := cmd.Flags()
	flags.StringSliceVarP(&listOpts.Filters, "filter", "f", []string{}, "Filter output based on conditions provided")
	flags.StringVar(&listOpts.Format, "format", "", "Pretty-print secrets using a Go template")
	flags.BoolVarP(&listOpts.Quiet, "quiet", "q", false, "Only display IDs")

	return cmd
}

func list(c cli.Interface, cmd *cobra.Command) error {
	conn := c.ClientConn()
	client := secret.NewSecretServiceClient(conn)
	request := &secret.ListSecretsRequest{}
	resp, err := client.ListSecrets(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
		return fmt.Errorf("Error listing secret: %s", err)
	}

	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "NAME\t")
	for _, s := range resp.Secrets {
		fmt.Fprintf(w, "%s\t\n", s.Spec.Annotations.Name)
	}
	w.Flush()
	return nil
}
