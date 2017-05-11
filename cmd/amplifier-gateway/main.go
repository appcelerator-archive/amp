package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/api/rpc/cluster"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/appcelerator/amp/api/rpc/storage"
	"github.com/appcelerator/amp/api/rpc/version"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	amplifierEndpoint = "amplifier" + configuration.DefaultPort
)

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	eheaders := []string{"Authorization", "X-Custom-header", "Grpc-Metadata-Amp.token"}
	w.Header().Set("Access-Controls-Expose-Headers", strings.Join(eheaders, ","))
	aheaders := []string{"Content-Type", "Accept", "Authorization", "Grpc-Metadata-Amp.token"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(aheaders, ","))
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

	if err := account.RegisterAccountHandlerFromEndpoint(ctx, mux, amplifierEndpoint, opts); err != nil {
		log.Fatal(err)
	}
	if err := cluster.RegisterClusterHandlerFromEndpoint(ctx, mux, amplifierEndpoint, opts); err != nil {
		log.Fatal(err)
	}
	if err := logs.RegisterLogsHandlerFromEndpoint(ctx, mux, amplifierEndpoint, opts); err != nil {
		log.Fatal(err)
	}
	if err := stack.RegisterStackHandlerFromEndpoint(ctx, mux, amplifierEndpoint, opts); err != nil {
		log.Fatal(err)
	}
	if err := stats.RegisterStatsHandlerFromEndpoint(ctx, mux, amplifierEndpoint, opts); err != nil {
		log.Fatal(err)
	}
	if err := service.RegisterServiceHandlerFromEndpoint(ctx, mux, amplifierEndpoint, opts); err != nil {
		log.Fatal(err)
	}
	if err := storage.RegisterStorageHandlerFromEndpoint(ctx, mux, amplifierEndpoint, opts); err != nil {
		log.Fatal(err)
	}
	if err := version.RegisterVersionHandlerFromEndpoint(ctx, mux, amplifierEndpoint, opts); err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(":80", allowCORS(mux))
	return
}
