package cluster

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/appcelerator/amp/pkg/docker"
	"github.com/docker/docker/api/types"

	"golang.org/x/net/context"
)

// Server is used to implement cluster.ClusterServer
type Server struct {
	Docker *docker.Docker
}

// Create implements cluster.Server
func (s *Server) Create(ctx context.Context, in *CreateRequest) (*CreateReply, error) {
	log.Println("[cluster] Create", in.String())

	log.Println(in.Name)
	log.Println(string(in.Compose))

	log.Println("[cluster] Success: created cluster")
	return &CreateReply{}, nil
}

// List implements cluster.Server
func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {
	log.Println("[cluster] List", in.String())

	log.Println("[cluster] Success: list")
	return &ListReply{}, nil
}

// Status implements cluster.Server
func (s *Server) Status(ctx context.Context, in *StatusRequest) (*StatusReply, error) {
	log.Println("[cluster] Status", in.String())

	log.Println("[cluster] Success: list")
	return &StatusReply{}, nil
}

// Update implements cluster.Server
func (s *Server) Update(ctx context.Context, in *UpdateRequest) (*UpdateReply, error) {
	log.Println("[cluster] Update", in.String())

	log.Println("[cluster] Success: list")
	return &UpdateReply{}, nil
}

// Remove implements cluster.Server
func (s *Server) Remove(ctx context.Context, in *RemoveRequest) (*RemoveReply, error) {
	log.Println("[cluster] Remove", in.String())

	log.Println("[cluster] Success: removed", in.Id)
	return &RemoveReply{}, nil
}

// NodeList get cluster node list
func (s *Server) NodeList(ctx context.Context, in *NodeListRequest) (*NodeListReply, error) {
	log.Println("[cluster] NodeList", in.String())

	list, err := s.Docker.GetClient().NodeList(ctx, types.NodeListOptions{})
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
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
