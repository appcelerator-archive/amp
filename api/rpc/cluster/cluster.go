package cluster

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/filters"
	"docker.io/go-docker/api/types/swarm"
	"github.com/appcelerator/amp/pkg/cloud"
	"github.com/appcelerator/amp/pkg/cloud/aws"
	"github.com/appcelerator/amp/pkg/docker"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is used to implement cluster.ClusterServer
type Server struct {
	Docker   *docker.Docker
	Provider cloud.Provider
	Region   string
}

const (
	CoreStackName = "amp"
)

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

	// TODO

	//return &ListReply{}, nil
	log.Infoln("[cluster] NotImplemented: list")
	return nil, status.Errorf(codes.Unimplemented, "Not implemented yet")
}

// swarmNodeStatus returns the swarm status for this node
func swarmNodeStatus(c *docker.Docker) (swarm.LocalNodeState, error) {
	info, err := c.GetClient().Info(context.Background())
	if err != nil {
		return "", err
	}
	return info.Swarm.LocalNodeState, nil
}

// infoAMPCore returns the number of AMP core services
func infoAMPCore(ctx context.Context, c *docker.Docker) (int, error) {
	var count int
	services, err := c.GetClient().ServiceList(ctx, types.ServiceListOptions{})
	if err != nil {
		return 0, err
	}
	for _, service := range services {
		if strings.HasPrefix(service.Spec.Name, fmt.Sprintf("%s_", CoreStackName)) {
			count++
		}
	}
	return count, nil
}

// infoUser returns the number of user services
func infoUser(ctx context.Context, c *docker.Docker) (int, error) {
	var count int
	services, err := c.GetClient().ServiceList(ctx, types.ServiceListOptions{})
	if err != nil {
		return 0, err
	}
	for _, service := range services {
		if !strings.HasPrefix(service.Spec.Name, CoreStackName) {
			count++
		}
	}
	return count, err
}

// Status implements cluster.Server
func (s *Server) ClusterStatus(ctx context.Context, in *StatusRequest) (*StatusReply, error) {
	log.Infoln("[cluster] Status", in.String())

	// Check the node status
	swarmStatus, err := swarmNodeStatus(s.Docker)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	// Assuming the swarm is not active
	coreServices := 0
	userServices := 0
	if swarmStatus == swarm.LocalNodeStateActive { // if it is, update the services
		ctx := context.Background()
		coreServices, err = infoAMPCore(ctx, s.Docker)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
		userServices, err = infoUser(ctx, s.Docker)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
	}
	stackInfo := map[string]string{"StackName": "local", "Region": s.Region, "NFSEndpoint": "disabled"}

	// check cloud providers
	switch s.Provider {
	case cloud.ProviderAWS:
		if err := aws.StackInfo(ctx, &stackInfo); err != nil {
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
	case cloud.ProviderLocal:
	default:
		return nil, status.Errorf(codes.Unimplemented, "provider [%s] is not yet implemented", string(s.Provider))
	}
	log.Infoln("[cluster] Success: status")
	return &StatusReply{
		Name:             stackInfo["StackName"],
		Provider:         string(s.Provider),
		Region:           s.Region,
		SwarmStatus:      string(swarmStatus),
		CoreServices:     strconv.Itoa(coreServices),
		UserServices:     strconv.Itoa(userServices),
		Endpoint:         stackInfo["DNSTarget"],
		NfsEndpoint:      stackInfo["NFSEndpoint"],
		InternalEndpoint: stackInfo["InternalDockerHost"],
		InternalPki:      stackInfo["InternalPKITarget"],
	}, nil
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

	filter := filters.NewArgs()
	if in.Id != "" {
		filter.Add("id", in.Id)
	}
	if in.Name != "" {
		filter.Add("name", in.Name)
	}
	if in.Role != "" {
		filter.Add("role", in.Role)
	}
	if in.EngineLabel != "" {
		filter.Add("label", in.EngineLabel)
	}
	list, err := s.Docker.GetClient().NodeList(ctx, types.NodeListOptions{Filters: filter})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	ret := &NodeListReply{}
	// prepare key value for label filtering
	lk := strings.Split(in.NodeLabel, "=")[0]
	lv := ""
	if strings.Contains(in.NodeLabel, "=") {
		lv = strings.Split(in.NodeLabel, "=")[1]
	}
	for _, node := range list {
		if lk != "" {
			labelMatch := false
			for k, v := range node.Spec.Annotations.Labels {
				if k == lk && (lv == "" || lv == v) {
					labelMatch = true
					break
				}
			}
			if !labelMatch {
				continue
			}
		}
		leader := false
		if node.ManagerStatus != nil {
			leader = node.ManagerStatus.Leader
		}
		var enginePlugins []*EnginePlugin
		for _, p := range node.Description.Engine.Plugins {
			enginePlugins = append(enginePlugins, &EnginePlugin{Type: p.Type, Name: p.Name})
		}
		ret.Nodes = append(ret.Nodes, &NodeReply{
			Id:            node.ID,
			Hostname:      node.Description.Hostname,
			Status:        string(node.Status.State),
			Availability:  string(node.Spec.Availability),
			Role:          string(node.Spec.Role),
			ManagerLeader: leader,
			NodeLabels:    node.Spec.Annotations.Labels,
			EngineLabels:  node.Description.Engine.Labels,
			NanoCpus:      node.Description.Resources.NanoCPUs,
			MemoryBytes:   node.Description.Resources.MemoryBytes,
			EngineVersion: node.Description.Engine.EngineVersion,
			EnginePlugins: enginePlugins,
		})
	}
	return ret, nil
}

// NodeCleanup removes nodes in the down state
func (s *Server) ClusterNodeCleanup(ctx context.Context, in *NodeCleanupRequest) (*NodeListReply, error) {
	log.Infoln("[cluster] NodeCleanup")

	list, err := s.Docker.GetClient().NodeList(ctx, types.NodeListOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	ret := &NodeListReply{}
	for _, node := range list {
		if node.Status.State == swarm.NodeStateDown {
			ret.Nodes = append(ret.Nodes, &NodeReply{
				Id:       node.ID,
				Hostname: node.Description.Hostname,
				Role:     string(node.Spec.Role),
			})
			// if the node is a manager, first demote it
			if node.Spec.Role == swarm.NodeRoleManager {
				log.Infoln("Demoting node", node.ID, node.Description.Hostname)
				node.Spec.Role = swarm.NodeRoleWorker
				if err = s.Docker.GetClient().NodeUpdate(ctx, node.ID, node.Version, node.Spec); err != nil {
					return nil, err
				}
			}
			if err = s.Docker.GetClient().NodeRemove(ctx, node.ID, types.NodeRemoveOptions{Force: in.Force}); err != nil {
				log.Infoln("Removing node", node.ID, node.Description.Hostname)
				return nil, err
			}
		}
	}
	return ret, nil
}
