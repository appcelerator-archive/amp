package agentcore

import (
	"fmt"
	"github.com/appcelerator/amp/cmd/adm-agent/agentgrpc"
	"github.com/docker/docker/api/types"
	//"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

//RegistryToken token used for registry
const RegistryToken = ""

// GetNodeInfo return node information
func (g *ClusterAgent) GetNodeInfo(ctx context.Context, req *agentgrpc.GetNodeInfoRequest) (*agentgrpc.NodeInfo, error) {
	return g.getNodeInfo(req)
}

// PurgeNode remove containers, images, volumes
func (g *ClusterAgent) PurgeNode(ctx context.Context, req *agentgrpc.PurgeNodeRequest) (*agentgrpc.PurgeNodeAnswer, error) {
	answer := &agentgrpc.PurgeNodeAnswer{}
	if req.Container {
		nb, err := g.purgeContainers(req.Force)
		if err != nil {
			return nil, err
		}
		answer.NbContainers = int32(nb)
	}
	if req.Volume {
		nb, err := g.purgeVolumes(req.Force)
		if err != nil {
			return nil, err
		}
		answer.NbVolumes = int32(nb)
	}
	if req.Image {
		nb, err := g.purgeImages(req.Force)
		if err != nil {
			return nil, err
		}
		answer.NbImages = int32(nb)
	}
	return answer, nil
}

// PullImage pull infra images
func (g *ClusterAgent) PullImage(ctx context.Context, req *agentgrpc.PullImageRequest) (*agentgrpc.AgentRet, error) {
	options := types.ImagePullOptions{}
	if RegistryToken != "" {
		options.RegistryAuth = RegistryToken
	}
	reader, err := g.dockerClient.ImagePull(g.ctx, req.Image, options)
	if err != nil {
		return nil, fmt.Errorf("image %s pull error: %v", req.Image, err)
	}
	data := make([]byte, 1000, 1000)
	for {
		_, err := reader.Read(data)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Errorf("Pull image %s error: %v", req.Image, err)
		}
	}
	return &agentgrpc.AgentRet{}, nil
}
