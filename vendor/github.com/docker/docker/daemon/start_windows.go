package daemon

import (
	"fmt"
	"path/filepath"

	"github.com/docker/docker/container"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/libcontainerd"
)

func (daemon *Daemon) getLibcontainerdCreateOptions(container *container.Container) (*[]libcontainerd.CreateOption, error) {
	createOptions := []libcontainerd.CreateOption{}

	// Are we going to run as a Hyper-V container?
	hvOpts := &libcontainerd.HyperVIsolationOption{}
	if container.HostConfig.Isolation.IsDefault() {
		// Container is set to use the default, so take the default from the daemon configuration
		hvOpts.IsHyperV = daemon.defaultIsolation.IsHyperV()
	} else {
		// Container is requesting an isolation mode. Honour it.
		hvOpts.IsHyperV = container.HostConfig.Isolation.IsHyperV()
	}

	// Generate the layer folder of the layer options
	layerOpts := &libcontainerd.LayerOption{}
	m, err := container.RWLayer.Metadata()
	if err != nil {
		return nil, fmt.Errorf("failed to get layer metadata - %s", err)
	}
	if hvOpts.IsHyperV {
		hvOpts.SandboxPath = filepath.Dir(m["dir"])
	} else {
		layerOpts.LayerFolderPath = m["dir"]
	}

	// Generate the layer paths of the layer options
	img, err := daemon.imageStore.Get(container.ImageID)
	if err != nil {
		return nil, fmt.Errorf("failed to graph.Get on ImageID %s - %s", container.ImageID, err)
	}
	// Get the layer path for each layer.
	max := len(img.RootFS.DiffIDs)
	for i := 1; i <= max; i++ {
		img.RootFS.DiffIDs = img.RootFS.DiffIDs[:i]
		layerPath, err := layer.GetLayerPath(daemon.layerStore, img.RootFS.ChainID())
		if err != nil {
			return nil, fmt.Errorf("failed to get layer path from graphdriver %s for ImageID %s - %s", daemon.layerStore, img.RootFS.ChainID(), err)
		}
		// Reverse order, expecting parent most first
		layerOpts.LayerPaths = append([]string{layerPath}, layerOpts.LayerPaths...)
	}

	// Now build the full set of options
	createOptions = append(createOptions, &libcontainerd.FlushOption{IgnoreFlushesDuringBoot: !container.HasBeenStartedBefore})
	createOptions = append(createOptions, hvOpts)
	createOptions = append(createOptions, layerOpts)

	return &createOptions, nil
}
