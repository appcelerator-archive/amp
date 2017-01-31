package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const port = ":3000"

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		hostname, err := os.Hostname()
		if err != nil {
			panic(err)
		}

		// Since HTTP/1.1 defaults to persistent connections, ensure we close the
		// connection with the response to make it easier to demo automatic
		// round-robin load balancing when refreshing in a browser
		// (this isn't an issue when using curl since it automatically closes
		// the connection).
		w.Header().Set("Connection", "close")
		response := fmt.Sprintf("[%s] pong", hostname)
		fmt.Fprintln(w, response)
	})

	fmt.Printf("listening on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
