package cluster

import (
	log "github.com/sirupsen/logrus"

	"docker.io/go-docker/api/types"
	"github.com/appcelerator/amp/pkg/docker"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is used to implement cluster.ClusterServer
type Server struct {
	Docker *docker.Docker
}

// Create implements cluster.Server
func (s *Server) ClusterCreate(ctx context.Context, in *CreateRequest) (*CreateReply, error) {
	log.Infoln("[cluster] Create", in.String())

	log.Infoln(in.Name)
	log.Infoln(string(in.Compose))

	log.Infoln("[cluster] Success: created cluster")
	return &CreateReply{}, nil
}

// List implements cluster.Server
func (s *Server) ClusterList(ctx context.Context, in *ListRequest) (*ListReply, error) {
	log.Infoln("[cluster] List", in.String())

	log.Infoln("[cluster] Success: list")
	return &ListReply{}, nil
}

// Status implements cluster.Server
func (s *Server) ClusterStatus(ctx context.Context, in *StatusRequest) (*StatusReply, error) {
	log.Infoln("[cluster] Status", in.String())

	log.Infoln("[cluster] Success: list")
	return &StatusReply{}, nil
}

// Update implements cluster.Server
func (s *Server) ClusterUpdate(ctx context.Context, in *UpdateRequest) (*UpdateReply, error) {
	log.Infoln("[cluster] Update", in.String())

	log.Infoln("[cluster] Success: list")
	return &UpdateReply{}, nil
}

// Remove implements cluster.Server
func (s *Server) ClusterRemove(ctx context.Context, in *RemoveRequest) (*RemoveReply, error) {
	log.Infoln("[cluster] Remove", in.String())

	log.Infoln("[cluster] Success: removed", in.Id)
	return &RemoveReply{}, nil
}

// NodeList get cluster node list
func (s *Server) ClusterNodeList(ctx context.Context, in *NodeListRequest) (*NodeListReply, error) {
	log.Infoln("[cluster] NodeList", in.String())

	list, err := s.Docker.GetClient().NodeList(ctx, types.NodeListOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	ret := &NodeListReply{}
	for _, node := range list {
		leader := false
		if node.ManagerStatus != nil {
			leader = node.ManagerStatus.Leader
		}
		ret.Nodes = append(ret.Nodes, &NodeReply{
			Id:            node.ID,
			Hostname:      node.Description.Hostname,
			Status:        string(node.Status.State),
			Availability:  string(node.Spec.Availability),
			Role:          string(node.Spec.Role),
			ManagerLeader: leader,
		})
	}
	return ret, nil
}
