package agentcore

import (
	"github.com/appcelerator/amp/cmd/adm-agent/agentgrpc"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func (g *ClusterAgent) getNodeInfo(req *agentgrpc.GetNodeInfoRequest) (*agentgrpc.NodeInfo, error) {
	ret := &agentgrpc.NodeInfo{}
	info, erri := g.dockerClient.Info(g.ctx)
	if erri != nil {
		return nil, erri
	}
	ret.Cpu = int64(info.NCPU)
	ret.Memory = info.MemTotal
	ret.NbContainers = int64(info.Containers)
	ret.NbContainersRunning = int64(info.ContainersRunning)
	ret.NbContainersPaused = int64(info.ContainersPaused)
	ret.NbContainersStopped = int64(info.ContainersStopped)
	ret.Images = int64(info.Images)
	return ret, nil
}

func (g *ClusterAgent) purgeContainers(force bool) (int, error) {
	nb := 0
	list, err := g.dockerClient.ContainerList(g.ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return 0, err
	}
	for _, cont := range list {
		if cont.State == "exited" {
			if err := g.dockerClient.ContainerRemove(g.ctx, cont.ID, types.ContainerRemoveOptions{Force: force}); err == nil {
				nb++
			} else {
				logf.error("Error removing container id=%s: %v\n", cont.ID, err)
			}

		}
	}
	return nb, nil
}

func (g *ClusterAgent) purgeVolumes(force bool) (int, error) {
	nb := 0
	filter := filters.NewArgs()
	list, err := g.dockerClient.VolumeList(g.ctx, filter)
	if err != nil {
		return 0, err
	}
	for _, vol := range list.Volumes {

		if vol.Name != "amp-registry" {
			if err := g.dockerClient.VolumeRemove(g.ctx, vol.Name, force); err == nil {
				nb++
			} else {
				logf.error("Error removing volume name=%s: %v\n", vol.Name, err)
			}

		}

	}
	return nb, nil
}

func (g *ClusterAgent) purgeImages(force bool) (int, error) {
	nb := 0
	list, err := g.dockerClient.ImageList(g.ctx, types.ImageListOptions{
		All: true,
	})
	if err != nil {
		return 0, err
	}
	for _, ima := range list {
		if _, err := g.dockerClient.ImageRemove(g.ctx, ima.ID, types.ImageRemoveOptions{Force: force, PruneChildren: true}); err == nil {
			nb++
		} else {
			logf.error("Error removing image name=%s: %v\n", ima.ID, err)
		}

	}
	return nb, nil
}
