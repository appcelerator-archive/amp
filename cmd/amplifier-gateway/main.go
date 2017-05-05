package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	amplifierEndpoint = "amplifier" + configuration.DefaultPort
)

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			// w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preFlightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preFlightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	return
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := stack.RegisterStackHandlerFromEndpoint(ctx, mux, amplifierEndpoint, opts)
	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(":80", allowCORS(mux))
	return
}
