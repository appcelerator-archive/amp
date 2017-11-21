package config

import (
	"context"
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/config"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
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

// NewListCommand returns a new instance of the list command for listing configs
func NewListCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Short:   "List configs",
		Aliases: []string{"list"},
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			//return list(c, cmd, listOpts)
			return list(c, cmd)
		},
	}

	flags := cmd.Flags()
	flags.StringSliceVarP(&listOpts.Filters, "filter", "f", []string{}, "Filter output based on conditions provided")
	flags.StringVar(&listOpts.Format, "format", "", "Pretty-print configs using a Go template")
	flags.BoolVarP(&listOpts.Quiet, "quiet", "q", false, "Only display IDs")

	return cmd
}

func list(c cli.Interface, cmd *cobra.Command) error {
	conn := c.ClientConn()
	client := config.NewConfigClient(conn)
	request := &config.ListRequest{}
	reply, err := client.ConfigList(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
		return fmt.Errorf("error listing config: %entry", err)
	}

	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\t")
	for _, entry := range reply.Entries {
		fmt.Fprintf(w, "%s\t%s\t\n", entry.Id, entry.Name)
	}
	w.Flush()
	return nil
}
