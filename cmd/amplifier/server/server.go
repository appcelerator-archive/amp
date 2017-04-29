package server

import (
	"fmt"
	"log"
	"net"
	"os"
	rt "runtime"
	"sync"

	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/api/rpc/cluster"
	"github.com/appcelerator/amp/api/rpc/function"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/appcelerator/amp/api/rpc/storage"
	"github.com/appcelerator/amp/api/rpc/version"
	"github.com/appcelerator/amp/api/runtime"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/functions"
	"github.com/appcelerator/amp/data/stacks"
	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/appcelerator/amp/pkg/docker"
	"github.com/appcelerator/amp/pkg/elasticsearch"
	"github.com/appcelerator/amp/pkg/mail"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"google.golang.org/grpc"
)

type (
	clientInitializer  func(*Configuration) error
	serviceInitializer func(*Configuration, *grpc.Server)
)

// Client initializers open connections to required backend services
// Clients are stored as members of runtime
var clientInitializers = []clientInitializer{
	initDocker,
	initElasticsearch,
	initEtcd,
	initMailer,
	initNats,
}

// Service initializers register the services with the grpc server
var serviceInitializers = []serviceInitializer{
	registerVersionServer,
	registerStorageServer,
	registerLogsServer,
	registerStackServer,
	registerStatsServer,
	registerFunctionServer,
	registerAccountServer,
	registerClusterServer,
}

// Start starts the amplifier server
func Start(c *Configuration) {
	// initialize clients
	initClients(c)

	// register services
	s := grpc.NewServer(
		grpc.StreamInterceptor(auth.StreamInterceptor),
		grpc.UnaryInterceptor(auth.Interceptor),
	)
	registerServices(c, s)

	// start listening
	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalf("Unable to listen on %s: %v\n", c.Port[1:], err)
	}
	log.Println("Listening on port:", c.Port[1:])
	log.Fatalln(s.Serve(lis))
}

func initClients(config *Configuration) {
	// ensure all initialization code fails fast on errors; there is no point in
	// attempting to continue in a degraded state if there are problems at start up

	var wg sync.WaitGroup
	wg.Add(len(clientInitializers))
	for _, f := range clientInitializers {
		go func(f clientInitializer) {
			defer wg.Done()
			if err := f(config); err != nil {
				log.Fatalln(err)
			}
		}(f)
	}

	// Wait for all inits to complete.
	wg.Wait()
}

func initEtcd(config *Configuration) error {
	runtime.Store = etcd.New(config.EtcdEndpoints, "amp", DefaultTimeout)
	return nil
}

func initElasticsearch(config *Configuration) error {
	runtime.Elasticsearch = elasticsearch.NewClient(config.ElasticsearchURL, DefaultTimeout)
	return nil
}

func initNats(config *Configuration) error {
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("unable to get hostname: %v", err)
	}
	runtime.NatsStreaming = ns.NewClient(config.NatsURL, ns.ClusterID, os.Args[0]+"-"+hostname, DefaultTimeout)
	return nil
}

func initDocker(config *Configuration) error {
	runtime.Docker = docker.NewClient(config.DockerURL, config.DockerVersion)
	log.Printf("Connecting to Docker API at %s version API: %s\n", config.DockerURL, config.DockerVersion)
	if err := runtime.Docker.Connect(); err != nil {
		return err
	}
	log.Println("Connected to Docker API at", config.DockerURL)
	return nil
}

func initMailer(config *Configuration) error {
	runtime.Mailer = mail.NewMailer(config.EmailKey, config.EmailSender)
	return nil
}

func registerServices(c *Configuration, s *grpc.Server) {
	var wg sync.WaitGroup
	wg.Add(len(serviceInitializers))
	for _, f := range serviceInitializers {
		go func(f serviceInitializer) {
			defer wg.Done()
			f(c, s)
		}(f)
	}

	// Wait for all service registrations to complete.
	wg.Wait()
}

func registerVersionServer(c *Configuration, s *grpc.Server) {
	version.RegisterVersionServer(s, &version.Server{
		Info: &version.Info{
			Version:   c.Version,
			Build:     c.Build,
			GoVersion: rt.Version(),
			Os:        rt.GOOS,
			Arch:      rt.GOARCH,
		},
	})
}

func registerLogsServer(c *Configuration, s *grpc.Server) {
	logs.RegisterLogsServer(s, &logs.Server{
		Es:            runtime.Elasticsearch,
		NatsStreaming: runtime.NatsStreaming,
	})
}

func registerStorageServer(c *Configuration, s *grpc.Server) {
	storage.RegisterStorageServer(s, &storage.Server{
		Store: runtime.Store,
	})
}

func registerStatsServer(c *Configuration, s *grpc.Server) {
	stats.RegisterStatsServer(s, &stats.Stats{
		Es: runtime.Elasticsearch,
	})
}

func registerFunctionServer(c *Configuration, s *grpc.Server) {
	function.RegisterFunctionServer(s, &function.Server{
		Functions:     functions.NewStore(runtime.Store),
		NatsStreaming: runtime.NatsStreaming,
	})
}

func registerAccountServer(c *Configuration, s *grpc.Server) {
	account.RegisterAccountServer(s, &account.Server{
		Accounts: accounts.NewStore(runtime.Store),
		Mailer:   runtime.Mailer,
	})
}

func registerStackServer(c *Configuration, s *grpc.Server) {
	stack.RegisterStackServer(s, &stack.Server{
		Accounts: accounts.NewStore(runtime.Store),
		Docker:   runtime.Docker,
		Stacks:   stacks.NewStore(runtime.Store),
	})
}

func registerClusterServer(c *Configuration, s *grpc.Server) {
	cluster.RegisterClusterServer(s, &cluster.Server{
		Docker: runtime.Docker,
	})
}
