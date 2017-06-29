package core

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
)

// ContainerData data
type ContainerData struct {
	name                     string
	ID                       string
	shortName                string
	serviceName              string
	serviceID                string
	stackName                string
	taskID                   string
	taskSlot                 int
	nodeID                   string
	role                     string
	pid                      int
	state                    string
	health                   string
	logsStream               io.ReadCloser
	logsReadError            bool
	metricsStream            io.ReadCloser
	metricsReadError         bool
	previousIOStats          *IOStats
	previousNetStats         *NetStats
	lastDateSaveTime         time.Time
	labels                   map[string]string
	squashedMetricsMessage   stats.MetricsEntry
	squashedMetricsMessageNb int64
}

// Verify if the event stream is working, if not start it
func (a *Agent) updateEventsStream() {
	if !a.eventStreamReading {
		log.Infoln("Opening docker events stream...")
		args := filters.NewArgs()
		args.Add("type", "container")
		args.Add("event", "die")
		args.Add("event", "stop")
		args.Add("event", "destroy")
		args.Add("event", "kill")
		args.Add("event", "create")
		args.Add("event", "start")
		eventsOptions := types.EventsOptions{Filters: args}
		stream, err := a.dock.GetClient().Events(context.Background(), eventsOptions)
		a.startEventStream(stream, err)
	}
}

// Start and read the docker event stream and update container list accordingly
func (a *Agent) startEventStream(stream <-chan events.Message, errs <-chan error) {
	a.eventStreamReading = true
	log.Infoln("start events stream reader")
	go func() {
		for {
			select {
			case err := <-errs:
				if err != nil {
					log.Errorf("Error reading event: %v\n", err)
					a.eventStreamReading = false
					return
				}
			case event := <-stream:
				log.Debugf("Docker event: action=%s containerId=%s\n", event.Action, event.Actor.ID)
				a.updateContainerMap(event.Action, event.Actor.ID)
			}
		}
	}()
}

// Update containers list considering event action and event container id
func (a *Agent) updateContainerMap(action string, containerID string) {
	if action == "start" {
		a.addContainer(containerID)
	} else if action == "destroy" || action == "die" || action == "kill" || action == "stop" {
		go func() {
			time.Sleep(5 * time.Second)
			a.removeContainer(containerID)
		}()
	}
}

// Add a container to the main container map and retrieve some container information
func (a *Agent) addContainer(ID string) {
	_, present := a.containers[ID]
	if present {
		return
	}
	container, err := a.dock.GetClient().ContainerInspect(context.Background(), ID)
	if err != nil {
		log.Errorf("Container inspect error: %v\n", err)
		return
	}

	// Create container data
	data := ContainerData{
		ID:            ID,
		name:          a.cleanName(container.Name),
		state:         container.State.Status,
		pid:           container.State.Pid,
		health:        "",
		logsStream:    nil,
		logsReadError: false,
	}

	data.squashedMetricsMessage.Cpu = &stats.MetricsCPUEntry{}
	data.squashedMetricsMessage.Mem = &stats.MetricsMemEntry{}
	data.squashedMetricsMessage.Net = &stats.MetricsNetEntry{}
	data.squashedMetricsMessage.Io = &stats.MetricsIOEntry{}
	a.clearMetricsMessage(&data)
	labels := container.Config.Labels
	data.serviceName = a.getMapValue(labels, "com.docker.swarm.service.name")
	//data.serviceName = strings.TrimPrefix(labels["com.docker.swarm.service.name"], labels["com.docker.stack.namespace"]+"_")
	if data.serviceName == "" {
		data.serviceName = "noService"
	}
	data.shortName = fmt.Sprintf("%s_%s", data.serviceName, ID[0:6])
	data.labels = labels
	data.serviceID = a.getMapValue(labels, "com.docker.swarm.service.id")
	data.taskID = a.getMapValue(labels, "com.docker.swarm.task.id")
	if data.taskID != "" {
		task, _, err := a.dock.GetClient().TaskInspectWithRaw(context.Background(), data.taskID)
		if err != nil {
			log.Errorf("Task inspect error: %v\n", err)
		} else {
			data.taskSlot = task.Slot
		}
	}
	data.nodeID = a.getMapValue(labels, "com.docker.swarm.node.id")
	data.stackName = a.getMapValue(labels, "com.docker.stack.namespace")
	if data.stackName == "" {
		data.stackName = "noStack"
	}
	data.role = a.getMapValue(labels, "io.amp.role")
	if container.State.Health != nil {
		data.health = container.State.Health.Status
	}
	if data.role == "infrastructure" {
		log.Infof("add infrastructure container %s\n", data.name)
	} else {
		log.Infof("add user container %s, stack=%s service=%s\n", data.name, data.stackName, data.serviceName)
	}
	data.labels = labels

	// Add the container data to the map
	a.containers[ID] = &data
}

// Strips '/' from beginning of container name, if present
func (a *Agent) cleanName(name string) string {
	if len(name) > 1 && name[0] == '/' {
		return name[1:]
	}
	return name
}

// Remove a container from the main container map
func (a *Agent) removeContainer(ID string) {
	data, ok := a.containers[ID]
	if ok {
		log.Infoln("Removing container", data.name)
		delete(a.containers, ID)
	}
	err := os.Remove(path.Join(containersDataDir, ID))
	if err != nil {
		log.Errorln("Error removing container", err)
	}
}

// Update container status and health
// TODO
// nolint: unused
func (a *Agent) updateContainer(id string) {
	data, ok := a.containers[id]
	if ok {
		inspect, err := a.dock.GetClient().ContainerInspect(context.Background(), id)
		if err == nil {
			// labels = inspect.Config.Labels
			data.state = inspect.State.Status
			data.health = ""
			if inspect.State.Health != nil {
				data.health = inspect.State.Health.Status
			}
			log.Infoln("Updating container", data.name)
		} else {
			log.Errorf("Container %s inspect error: %v\n", data.name, err)
		}
	}
}

func (a *Agent) getMapValue(labelMap map[string]string, name string) string {
	if val, exist := labelMap[name]; exist {
		return val
	}
	return ""
}
