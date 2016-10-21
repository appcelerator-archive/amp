package core

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"golang.org/x/net/context"
)

//Verify if the event stream is working, if not start it
func updateEventsStream() {
	if !agent.eventStreamReading {
		fmt.Println("Opening docker events stream...")
		args := filters.NewArgs()
		args.Add("type", "container")
		args.Add("event", "die")
		args.Add("event", "stop")
		args.Add("event", "destroy")
		args.Add("event", "kill")
		args.Add("event", "create")
		args.Add("event", "start")
		eventsOptions := types.EventsOptions{Filters: args}
		stream, err := agent.dockerClient.Events(context.Background(), eventsOptions)
		startEventStream(stream, err)
	}
}

// Start and read the docker event stream and update container list accordingly
func startEventStream(stream <-chan events.Message, errs <-chan error) {
	agent.eventStreamReading = true
	fmt.Println("start events stream reader")
	go func() {
		for {
			select {
			case err := <-errs:
				if err != nil {
					fmt.Printf("Error reading event: %v\n", err)
					agent.eventStreamReading = false
					return
				}
			case event := <-stream:
				fmt.Printf("Docker event: action=%s containerId=%s\n", event.Action, event.Actor.ID)
				agent.updateContainerMap(event.Action, event.Actor.ID)
			}
		}
	}()
}
