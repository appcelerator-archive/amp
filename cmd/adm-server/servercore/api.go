package servercore

import (
	"fmt"
	"net/http"
)

const baseURL = "/api/v1"

//Start API server
func initAPI() {
	fmt.Println("Start API server on port " + conf.apiPort)
	go func() {
		http.HandleFunc(baseURL+"/health", serverHealth)
		http.ListenAndServe(":"+conf.apiPort, nil)
	}()
}

//for HEALTHCHECK Dockerfile instruction
//for HEALTHCHECK Dockerfile instruction
func serverHealth(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(200)
	//resp.WriteHeader(400)
}
