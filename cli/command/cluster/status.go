package cluster

import (
	"context"
	"errors"
	"fmt"

	"github.com/appcelerator/amp/api/rpc/cluster"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/pkg/cloud"
	"github.com/spf13/cobra"
	grpcStatus "google.golang.org/grpc/status"
)

// NewStatusCommand returns a new instance of the status command for querying the state of amp cluster.
func NewStatusCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Aliases: []string{"info"},
		Short:   "Retrieve details about an amp cluster",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return status(c, cmd)
		},
	}

	return cmd
}

func status(c cli.Interface, cmd *cobra.Command) error {
	req := &cluster.StatusRequest{}
	client := cluster.NewClusterClient(c.ClientConn())
	reply, err := client.ClusterStatus(context.Background(), req)
	if err != nil {
		if s, ok := grpcStatus.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	if reply.Provider == "" {
		return errors.New("empty reply from server, probably an API mismatch")
	}
	fmt.Printf("Provider:      %s\n", reply.Provider)
	if reply.Provider != string(cloud.ProviderLocal) {
		fmt.Printf("Cluster name:  %s\n", reply.Name)
		fmt.Printf("Region:        %s\n", reply.Region)
	}
	fmt.Printf("Swarm Status:  %s\n", reply.SwarmStatus)
	fmt.Printf("Core Services: %s\n", reply.CoreServices)
	fmt.Printf("User Services: %s\n", reply.UserServices)
	if reply.Endpoint != "" {
		fmt.Printf("DNS Target:    %s\n", reply.Endpoint)
	}
	if reply.NfsEndpoint != "disabled" {
		fmt.Printf("NFS Endpoint:  %s\n", reply.NfsEndpoint)
	}
	return nil
}
