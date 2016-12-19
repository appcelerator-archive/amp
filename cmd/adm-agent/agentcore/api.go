package agentcore

import (
	"fmt"
	"net/http"
)

const baseURL = "/api/v1"

type agentAPI struct {
	agent *ClusterAgent
}

//Start API server
func (g *agentAPI) initAPI(agent *ClusterAgent) {
	g.agent = agent
	fmt.Println("Start API server on port " + conf.apiPort)
	go func() {
		http.HandleFunc(baseURL+"/health", g.agentHealth)
		http.ListenAndServe(":"+conf.apiPort, nil)
	}()
}

//for HEALTHCHECK Dockerfile instruction
func (g *agentAPI) agentHealth(resp http.ResponseWriter, req *http.Request) {
	if g.agent.healthy {
		resp.WriteHeader(200)
	} else {
		resp.WriteHeader(400)
	}
}
