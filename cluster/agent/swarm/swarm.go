package swarmkit

import (
	"context"
	"crypto/tls"
	"net"
	"os"
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/cmd/swarmd/defaults"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type NodeMap map[string]*api.Node

type NodeFilter func(n *api.Node) bool

func AllNodesFilter(n *api.Node) bool {
	return n != nil
}

func LiveNodesFilter(n *api.Node) bool {
	return n != nil && n.Status.State != api.NodeStatus_DOWN
}

const WatchActionKindAll =
	api.WatchActionKindUnknown |
	api.WatchActionKindCreate |
	api.WatchActionKindUpdate |
	api.WatchActionKindRemove


// /var/run/docker/swarm/control.sock
func DefaultSocket() string {
	swarmSocket := os.Getenv("SWARM_SOCKET")
	if swarmSocket != "" {
		return swarmSocket
	}
	return defaults.ControlAPISocket
}

// Dial establishes a connection and creates a client.
func Dial(addr string) (api.ControlClient, *grpc.ClientConn, error) {
	conn, err := DialConn(addr)
	if err != nil {
		return nil, nil, err
	}

	return api.NewControlClient(conn), conn, nil
}

// DialConn establishes a connection to SwarmKit.
func DialConn(addr string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{}
	insecureCreds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	opts = append(opts, grpc.WithTransportCredentials(insecureCreds))
	opts = append(opts, grpc.WithDialer(
		func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}))
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func Context() context.Context {
	// TODO: create a more appropriate context
	return context.TODO()
}

func ListNodes(ctx context.Context, c api.ControlClient, filter NodeFilter) (NodeMap, error) {
	if filter == nil {
		filter = AllNodesFilter
	}

	nr, err := c.ListNodes(ctx, &api.ListNodesRequest{})
	if err != nil {
		return nil, err
	}

	nodes := make(NodeMap)
	for _, n := range nr.Nodes {
		if filter(n) {
			nodes[n.ID] = n
		}
	}

	return nodes, nil
}

func NewWatchRequest(entries []*api.WatchRequest_WatchEntry, resumeFrom *api.Version, includeOldObject bool) *api.WatchRequest {
	return &api.WatchRequest{
		Entries: entries,
		ResumeFrom: resumeFrom,
		IncludeOldObject: includeOldObject,
	}
}


// @param kind
// node
// service
// network
// task
// cluster
// secret
// resource
// extension
// config
//
// @param action
// WatchActionKindUnknown WatchActionKind = 0
// WatchActionKindCreate  WatchActionKind = 1
// WatchActionKindUpdate  WatchActionKind = 2
// WatchActionKindRemove  WatchActionKind = 4
//
// @param filters
// SelectBy { By: }
// Types that are valid to be assigned to SelectBy.By
//	*SelectBy_ID
//	*SelectBy_IDPrefix
//	*SelectBy_Name
//	*SelectBy_NamePrefix
//	*SelectBy_Custom
//	*SelectBy_CustomPrefix
//	*SelectBy_ServiceID
//	*SelectBy_NodeID
//	*SelectBy_Slot
//	*SelectBy_DesiredState
//	*SelectBy_Role
//	*SelectBy_Membership
//	*SelectBy_ReferencedNetworkID
//	*SelectBy_ReferencedSecretID
//	*SelectBy_ReferencedConfigID
//	*SelectBy_Kind
//
func NewWatchRequestEntry(kind string, action api.WatchActionKind, filters []*api.SelectBy) *api.WatchRequest_WatchEntry {
	return &api.WatchRequest_WatchEntry{
		Kind: kind,
		Action: action,
		Filters: filters,
	}
}

