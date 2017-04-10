package core

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
)

// ContainerData data
type ContainerData struct {
	name             string
	ID               string
	shortName        string
	serviceName      string
	serviceID        string
	stackName        string
	taskID           string
	nodeID           string
	role             string
	pid              int
	state            string
	health           string
	logsStream       io.ReadCloser
	logsReadError    bool
	metricsStream    io.ReadCloser
	metricsReadError bool
	previousIOStats  *IOStats
	previousNetStats *NetStats
	lastDateSaveTime time.Time
}

// Verify if the event stream is working, if not start it
func (a *Agent) updateEventsStream() {
	if !a.eventStreamReading {
		log.Println("Opening docker events stream...")
		args := filters.NewArgs()
		args.Add("type", "container")
		args.Add("event", "die")
		args.Add("event", "stop")
		args.Add("event", "destroy")
		args.Add("event", "kill")
		args.Add("event", "create")
		args.Add("event", "start")
		eventsOptions := types.EventsOptions{Filters: args}
		stream, err := a.dockerClient.Events(context.Background(), eventsOptions)
		a.startEventStream(stream, err)
	}
}

// Start and read the docker event stream and update container list accordingly
func (a *Agent) startEventStream(stream <-chan events.Message, errs <-chan error) {
	a.eventStreamReading = true
	log.Println("start events stream reader")
	go func() {
		for {
			select {
			case err := <-errs:
				if err != nil {
					log.Printf("Error reading event: %v\n", err)
					a.eventStreamReading = false
					return
				}
			case event := <-stream:
				log.Printf("Docker event: action=%s containerId=%s\n", event.Action, event.Actor.ID)
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
	_, ok := a.containers[ID]
	if !ok {
		inspect, err := a.dockerClient.ContainerInspect(context.Background(), ID)
		if err == nil {
			data := ContainerData{
				ID:            ID,
				name:          a.cleanName(inspect.Name),
				state:         inspect.State.Status,
				pid:           inspect.State.Pid,
				health:        "",
				logsStream:    nil,
				logsReadError: false,
			}
			labels := inspect.Config.Labels
			// data.serviceName = a.getMapValue(labels, "com.docker.swarm.service.name")
			data.serviceName = strings.TrimPrefix(labels["com.docker.swarm.service.name"], labels["com.docker.stack.namespace"]+"_")
			if data.serviceName == "" {
				data.serviceName = "noService"
			}
			data.shortName = fmt.Sprintf("%s_%s", data.serviceName, ID[0:6])
			data.serviceID = a.getMapValue(labels, "com.docker.swarm.service.id")
			data.taskID = a.getMapValue(labels, "com.docker.swarm.task.id")
			data.nodeID = a.getMapValue(labels, "com.docker.swarm.node.id")
			data.stackName = a.getMapValue(labels, "com.docker.stack.namespace")
			if data.stackName == "" {
				data.stackName = "noStack"
			}
			data.role = a.getMapValue(labels, "io.amp.role")
			if inspect.State.Health != nil {
				data.health = inspect.State.Health.Status
			}
			if data.role == "infrastructure" {
				log.Printf("add infrastructure container  %s\n", data.name)
			} else {
				log.Printf("add user container %s, stack=%s service=%s\n", data.name, data.stackName, data.serviceName)
			}
			a.containers[ID] = &data
		} else {
			log.Printf("Container inspect error: %v\n", err)
		}
	}
}

//Remove the charactere '/' form the beginning of container name if exist
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
		log.Println("remove container", data.name)
		delete(a.containers, ID)
	}
	err := os.Remove(path.Join(containersDataDir, ID))
	if err != nil {
		log.Println("Error removing container data directory: ", err)
	}
}

// Update container status and health
// TODO
// nolint: unused
func (a *Agent) updateContainer(id string) {
	data, ok := a.containers[id]
	if ok {
		inspect, err := a.dockerClient.ContainerInspect(context.Background(), id)
		if err == nil {
			// labels = inspect.Config.Labels
			data.state = inspect.State.Status
			data.health = ""
			if inspect.State.Health != nil {
				data.health = inspect.State.Health.Status
			}
			log.Println("update container", data.name)
		} else {
			log.Printf("Container %s inspect error: %v\n", data.name, err)
		}
	}
}

func (a *Agent) getMapValue(labelMap map[string]string, name string) string {
	if val, exist := labelMap[name]; exist {
		return val
	}
	return ""
}
