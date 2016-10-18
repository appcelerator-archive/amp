package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type apiContainer struct {
	ContainerID string
	ServiceName string
	ServiceID   string
	State       string
	Health      string
}

const baseURL = "/api/v1"

//Start API server
func initAPI() {
	fmt.Println("Start API server on port " + conf.apiPort)
	go func() {
		http.HandleFunc(baseURL+"/health", agentHealth)
		http.HandleFunc(baseURL+"/containers", getHandledContainers)
		http.ListenAndServe(":"+conf.apiPort, nil)
	}()
}

//for HEALTHCHECK Dockerfile instruction
func agentHealth(resp http.ResponseWriter, req *http.Request) {
	if agent.eventStreamReading {
		resp.WriteHeader(200)
	} else {
		fmt.Println("execute /health: return not healthy")
		resp.WriteHeader(400)
	}
}

//return the running container list with their paremeter including health
func getHandledContainers(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("execute api /api/v1/containers")
	containers := make([]apiContainer, len(agent.containers))
	var nn int
	if time.Since(agent.lastUpdate) > time.Duration(3)*time.Second {
		for key := range agent.containers {
			agent.updateContainer(key)
		}
		agent.lastUpdate = time.Now()
	}
	for key, data := range agent.containers {
		containers[nn] = apiContainer{
			ContainerID: key,
			ServiceName: data.labels["com.docker.swarm.service.name"],
			ServiceID:   data.labels["com.docker.swarm.service.id"],
			State:       data.state,
			Health:      data.health,
		}
		nn++
	}
	json.NewEncoder(resp).Encode(containers)
}
