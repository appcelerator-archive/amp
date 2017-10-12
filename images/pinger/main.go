package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const port = ":3000"

var (
	httpResponsesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "pinger",
			Subsystem: "server",
			Name:      "http_responses_total",
			Help:      "The count of http responses issued",
		},
	)
	httpResponseLatencies = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "pinger",
			Subsystem: "server",
			Name:      "http_response_latencies",
			Help:      "The latency of http responses issued",
		},
	)
)

func ping(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
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
	elapsed := time.Since(start)
	msElapsed := elapsed / time.Millisecond
	httpResponsesTotal.Inc()
	httpResponseLatencies.Observe(float64(msElapsed))
}

func main() {
	prometheus.MustRegister(httpResponsesTotal)
	prometheus.MustRegister(httpResponseLatencies)
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/ping", ping)
	fmt.Printf("listening on %s\n", port)

	log.Fatal(http.ListenAndServe(port, nil))
}
