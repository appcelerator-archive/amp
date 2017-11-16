package docker

import (
	"strings"

	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/swarm"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/manager/constraint"
	"golang.org/x/net/context"
)

// NodeInspect inspects a node
func (d *Docker) NodeInspect(ctx context.Context, nodeID string) (swarm.Node, error) {
	nodeEntity, _, err := d.client.NodeInspectWithRaw(ctx, nodeID)
	if err != nil {
		return swarm.Node{}, err
	}
	return nodeEntity, nil
}

// NodeList list the nodes
func (d *Docker) NodeList(ctx context.Context, options types.NodeListOptions) ([]swarm.Node, error) {
	return d.client.NodeList(ctx, options)
}

// ExpectedNumberOfTasks returns expected number of tasks of a service
func (d *Docker) ExpectedNumberOfTasks(ctx context.Context, serviceID string) (int32, error) {
	var expectedTasks int32
	serviceInfo, err := d.ServiceInspect(ctx, serviceID)
	if err != nil {
		return 0, err
	}
	matchingNodeCount, err := d.numberOfMatchingNodes(ctx, serviceInfo)
	if err != nil {
		return 0, err
	}
	if matchingNodeCount == 0 {
		return 0, nil
	}
	if serviceInfo.Spec.Mode.Global != nil {
		expectedTasks = matchingNodeCount
	} else {
		expectedTasks = int32(*serviceInfo.Spec.Mode.Replicated.Replicas)
	}
	return expectedTasks, nil
}

// numberOfMatchingNodes returns number of nodes matching placement constraints
func (d *Docker) numberOfMatchingNodes(ctx context.Context, serviceInfo swarm.Service) (int32, error) {
	// placement constraints
	constraints, _ := constraint.Parse(serviceInfo.Spec.TaskTemplate.Placement.Constraints)
	// list all nodes in the swarm
	nodes, err := d.client.NodeList(ctx, types.NodeListOptions{})
	if err != nil {
		return 0, err
	}
	// inspect every node on the swarm to check for satisfying constraints
	var matchingNodes int32
	for _, node := range nodes {
		apiNode := d.nodeToNode(ctx, node)
		if constraint.NodeMatches(constraints, apiNode) {
			matchingNodes++
		}
	}
	return matchingNodes, nil
}

// nodeToNode converts a swarm node type into an api node type
func (d *Docker) nodeToNode(ctx context.Context, swarmNode swarm.Node) *api.Node {
	apiNode := &api.Node{
		ID: swarmNode.ID,
		Status: api.NodeStatus{
			Addr:  swarmNode.Status.Addr,
			State: api.NodeStatus_State(api.NodeStatus_State_value[strings.ToUpper(string(swarmNode.Status.State))]),
		},
		Spec: api.NodeSpec{
			Availability: api.NodeSpec_Availability(api.NodeSpec_Availability_value[strings.ToUpper(string(swarmNode.Spec.Availability))]),
			Annotations: api.Annotations{
				Labels: swarmNode.Spec.Labels,
			},
		},
		Description: &api.NodeDescription{
			Hostname: swarmNode.Description.Hostname,
			Platform: &api.Platform{
				OS:           swarmNode.Description.Platform.OS,
				Architecture: swarmNode.Description.Platform.Architecture,
			},
			Engine: &api.EngineDescription{
				EngineVersion: swarmNode.Description.Engine.EngineVersion,
				Labels:        swarmNode.Description.Engine.Labels,
			},
		},
		Role: api.NodeRole(api.NodeRole_value[strings.ToUpper(string(swarmNode.Spec.Role))]),
	}
	if swarmNode.ManagerStatus != nil {
		apiNode.ManagerStatus = &api.ManagerStatus{
			Leader:       swarmNode.ManagerStatus.Leader,
			Addr:         swarmNode.ManagerStatus.Addr,
			Reachability: api.RaftMemberStatus_Reachability(api.RaftMemberStatus_Reachability_value[strings.ToUpper(string(swarmNode.ManagerStatus.Reachability))]),
		}
	}

	for _, plugin := range swarmNode.Description.Engine.Plugins {
		apiNode.Description.Engine.Plugins = append(apiNode.Description.Engine.Plugins, api.PluginDescription{Type: plugin.Type, Name: plugin.Name})
	}
	return apiNode
}
