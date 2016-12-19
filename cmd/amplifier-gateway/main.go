package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/appcelerator/amp/api/rpc/function"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/appcelerator/amp/api/rpc/topic"
)

var (
	amplifierEndpoint = flag.String("amplifier_endpoint", "localhost:8080", "endpoint of amplifier")
)

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			// amp haproxy is doing some CORS headers as well
			// w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	return
}

func run() (err error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err = logs.RegisterLogsHandlerFromEndpoint(ctx, mux, *amplifierEndpoint, opts)
	if err != nil {
		return
	}
	err = service.RegisterServiceHandlerFromEndpoint(ctx, mux, *amplifierEndpoint, opts)
	if err != nil {
		return
	}
	err = stack.RegisterStackServiceHandlerFromEndpoint(ctx, mux, *amplifierEndpoint, opts)
	if err != nil {
		return
	}
	err = stats.RegisterStatsHandlerFromEndpoint(ctx, mux, *amplifierEndpoint, opts)
	if err != nil {
		return
	}
	err = topic.RegisterTopicHandlerFromEndpoint(ctx, mux, *amplifierEndpoint, opts)
	if err != nil {
		return
	}
	err = function.RegisterFunctionHandlerFromEndpoint(ctx, mux, *amplifierEndpoint, opts)
	if err != nil {
		return
	}

	http.ListenAndServe(":3000", allowCORS(mux))
	return
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
