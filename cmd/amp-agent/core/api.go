package core

import (
	"fmt"
	"net/http"
)

const baseURL = "/api/v1"

//Start API server
func (a *Agent) initAPI() {
	fmt.Println("Start API server on port " + conf.apiPort)
	go func() {
		http.HandleFunc(baseURL+"/health", a.agentHealth)
		http.ListenAndServe(":"+conf.apiPort, nil)
	}()
}

//for HEALTHCHECK Dockerfile instruction
func (a *Agent) agentHealth(resp http.ResponseWriter, req *http.Request) {
	if a.eventStreamReading {
		resp.WriteHeader(200)
	} else {
		fmt.Println("execute /health: return not healthy")
		resp.WriteHeader(400)
	}
}
