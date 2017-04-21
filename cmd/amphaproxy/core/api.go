package core

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

const port = "8090"

//Start API server
func initAPI() {
	fmt.Println("Start API server on port " + port)
	go func() {
		http.HandleFunc("/", receivedURL)
		http.ListenAndServe(":"+port, nil)
	}()
}

//for HEALTHCHECK Dockerfile instruction
func receivedURL(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(404)
	// is this an IP or a hostname?
	matched, err := regexp.MatchString(`([:digit:]+\.){4}`, req.Host)
	if err != nil {
		fmt.Printf("receivedURL: Failed to parse hostname: %v\n", err)
		return
	}
	if matched {
		fmt.Fprintln(resp, "Sorry, you can't access this service through an IP, please use a FQDN")
		return
	}
	list := strings.Split(req.Host, ".")
	if len(list) >= 2 {
		fmt.Fprintf(resp, "no server found for stack=%s service=%s\n", list[1], list[0])
		return
	}
	fmt.Fprintf(resp, "no server found for host: %s\n", req.Host)
}
