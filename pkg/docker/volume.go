package docker

import (
	"fmt"
	"time"

	"docker.io/go-docker/api/types"
	"github.com/appcelerator/amp/docker/cli/opts"
	"golang.org/x/net/context"
)

func (d *Docker) ListVolumes(filter opts.FilterOpt) ([]*types.Volume, error) {
	if d.GetClient() == nil {
		return nil, fmt.Errorf("Docker client is not connected")
	}
	reply, err := d.client.VolumeList(context.Background(), filter.Value())
	if err != nil {
		return nil, err
	}
	return reply.Volumes, nil
}

func (d *Docker) RemoveVolume(name string, force bool, retries int) error {
	success := false
	for i := 0; i < retries; i++ {
		if err := d.client.VolumeRemove(context.Background(), name, force); err != nil {
			time.Sleep(time.Second)
			continue
		}
		success = true
		break
	}
	if !success {
		return fmt.Errorf("timed out trying to remove volume: %s", name)
	}
	return nil
}
