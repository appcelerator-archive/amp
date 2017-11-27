package cluster

import (
	"context"
	"errors"
	"fmt"

	"github.com/appcelerator/amp/api/rpc/cluster"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	grpcStatus "google.golang.org/grpc/status"
)

// NewListCommand returns a new instance of the list command for amp clusters.
func NewListCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List deployed amp clusters",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(c)
		},
	}
}

func list(c cli.Interface) error {
	req := &cluster.ListRequest{}
	client := cluster.NewClusterClient(c.ClientConn())
	reply, err := client.ClusterList(context.Background(), req)
	if err != nil {
		if s, ok := grpcStatus.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	fmt.Printf("%+v", reply)
	return nil
}
