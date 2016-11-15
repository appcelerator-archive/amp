package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	//  "google.golang.org/grpc"
	//  "github.com/appcelerator/amp/api/rpc/logs"
)

var (
	amplifierEndpoint = flag.String("amplifier_endpoint", "localhost:8080", "endpoint of amplifier")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	//  opts := []grpc.DialOption{grpc.WithInsecure()}
	//  err := logs.RegisterLogsHandlerFromEndpoint(ctx, mux, *amplifierEndpoint, opts)
	//  if err != nil {
	//    return err
	//  }

	http.ListenAndServe(":3000", mux)
	return nil
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
