package image

import (
	"context"
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/image"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

// NewListCommand returns a new instance of the list command for listing images
func NewListCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Short:   "List repositories on the cluster registry",
		Aliases: []string{"list"},
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(c, cmd)
		},
	}
}

func list(c cli.Interface, cmd *cobra.Command) error {
	conn := c.ClientConn()
	client := image.NewImageClient(conn)
	request := &image.ListRequest{}
	reply, err := client.ImageList(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New("ImageList API call failed: " + s.Message())
		}
		return fmt.Errorf("error listing images: %s", err)
	}

	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "NAME\tTAG\tDIGEST")
	for _, repo := range reply.Entries {
		for _, i := range repo.Entries {
			fmt.Fprintf(w, "%s\t%s\t%s\n", repo.Name, i.Tag, i.Digest)
		}
	}
	w.Flush()
	return nil
}
