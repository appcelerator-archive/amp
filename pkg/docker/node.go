package docker

import (
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/swarm"
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
	constraints, _ := Parse(serviceInfo.Spec.TaskTemplate.Placement.Constraints)
	// list all nodes in the swarm
	nodes, err := d.client.NodeList(ctx, types.NodeListOptions{})
	if err != nil {
		return 0, err
	}
	// inspect every node on the swarm to check for satisfying constraints
	var matchingNodes int32
	for _, node := range nodes {
		if NodeMatches(constraints, &node) {
			matchingNodes++
		}
	}
	return matchingNodes, nil
}
