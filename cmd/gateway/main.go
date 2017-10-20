package main

import (
	"crypto/tls"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/api/rpc/cluster"
	"github.com/appcelerator/amp/api/rpc/dashboard"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/node"
	"github.com/appcelerator/amp/api/rpc/resource"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/rpc/version"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

const (
	listenAddress     = ":80"
	amplifierEndpoint = "amplifier" + configuration.DefaultPort
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the info level or above.
	log.SetLevel(log.InfoLevel)
}

// allowCORS allows Cross Origin Resource Sharing from any origin.
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

	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	creds := credentials.NewTLS(tlsConfig)
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Minute),
		grpc.WithCompressor(grpc.NewGZIPCompressor()),
		grpc.WithDecompressor(grpc.NewGZIPDecompressor()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                5 * time.Minute,
			Timeout:             20 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	if err := account.RegisterAccountHandlerFromEndpoint(ctx, mux, amplifierEndpoint, opts); err != nil {
		log.Fatal(err)
	}
	if err := cluster.RegisterClusterHandlerFromEndpoint(ctx, mux, amplifierEndpoint, opts); err != nil {
		log.Fatal(err)
	}
	if err := dashboard.RegisterDashboardHandlerFromEndpoint(ctx, mux, amplifierEndpoint, opts); err != nil {
		log.Fatal(err)
	}
	if err := logs.RegisterLogsHandlerFromEndpoint(ctx, mux, amplifierEndpoint, opts); err != nil {
		log.Fatal(err)
	}
	if err := resource.RegisterResourceHandlerFromEndpoint(ctx, mux, amplifierEndpoint, opts); err != nil {
		log.Fatal(err)
	}
	if err := stack.RegisterStackHandlerFromEndpoint(ctx, mux, amplifierEndpoint, opts); err != nil {
		log.Fatal(err)
	}
	if err := service.RegisterServiceHandlerFromEndpoint(ctx, mux, amplifierEndpoint, opts); err != nil {
		log.Fatal(err)
	}
	if err := node.RegisterNodeHandlerFromEndpoint(ctx, mux, amplifierEndpoint, opts); err != nil {
		log.Fatal(err)
	}
	if err := version.RegisterVersionHandlerFromEndpoint(ctx, mux, amplifierEndpoint, opts); err != nil {
		log.Fatal(err)
	}

	log.Infoln("gateway successfully initialized. Start listening on:", listenAddress)
	log.Fatalln(http.ListenAndServe(listenAddress, handlers.CompressHandler(handlers.CombinedLoggingHandler(os.Stdout, allowCORS(mux)))))
	return
}
