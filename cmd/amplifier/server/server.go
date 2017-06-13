package server

import (
	"log"
	"net"
	"os"
	"runtime"
	"sync"

	"fmt"

	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/api/rpc/cluster"
	"github.com/appcelerator/amp/api/rpc/dashboard"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/node"
	"github.com/appcelerator/amp/api/rpc/resource"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/appcelerator/amp/api/rpc/version"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/dashboards"
	"github.com/appcelerator/amp/data/stacks"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/appcelerator/amp/pkg/docker"
	"github.com/appcelerator/amp/pkg/elasticsearch"
	"github.com/appcelerator/amp/pkg/mail"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"google.golang.org/grpc"
)

type serviceInitializer func(*Amplifier, *grpc.Server)

// Amplifier represents the AMP gRPC server
type Amplifier struct {
	config     *configuration.Configuration
	docker     *docker.Docker
	storage    storage.Interface
	es         *elasticsearch.Elasticsearch
	ns         *ns.NatsStreaming
	mailer     *mail.Mailer
	tokens     *auth.Tokens
	accounts   accounts.Interface
	stacks     stacks.Interface
	dashboards dashboards.Interface
}

// Service initializers register the services with the grpc server
var serviceInitializers = []serviceInitializer{
	registerVersionServer,
	registerLogsServer,
	registerStackServer,
	registerStatsServer,
	registerAccountServer,
	registerClusterServer,
	registerServiceServer,
	registerNodeServer,
	registerResourceServer,
	registerDashboardServer,
}

func New(config *configuration.Configuration) (*Amplifier, error) {
	if config.JWTSecretKey == "" {
		return nil, fmt.Errorf("JWTSecret key cannot be empty. Please check amplifier configuration.")
	}
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	etcd := etcd.New(config.EtcdEndpoints, "amp", configuration.DefaultTimeout)
	accounts, err := accounts.NewStore(etcd, config.Registration, config.SUPassword)
	if err != nil {
		return nil, err
	}
	amp := &Amplifier{
		config:     config,
		storage:    etcd,
		es:         elasticsearch.NewClient(config.ElasticsearchURL, configuration.DefaultTimeout),
		ns:         ns.NewClient(config.NatsURL, ns.ClusterID, "amplifier-"+hostname, configuration.DefaultTimeout),
		docker:     docker.NewClient(config.DockerURL, config.DockerVersion),
		mailer:     mail.NewMailer(config.EmailKey, config.EmailSender, config.Notifications),
		tokens:     auth.New(config.JWTSecretKey),
		accounts:   accounts,
		stacks:     stacks.NewStore(etcd, accounts),
		dashboards: dashboards.NewStore(etcd, accounts),
	}
	return amp, nil
}

// Start starts the amplifier server
func (a *Amplifier) Start() {
	interceptors := &auth.Interceptors{Tokens: a.tokens}
	// register services
	s := grpc.NewServer(
		grpc.StreamInterceptor(interceptors.StreamInterceptor),
		grpc.UnaryInterceptor(interceptors.Interceptor),
		grpc.RPCCompressor(grpc.NewGZIPCompressor()),
		grpc.RPCDecompressor(grpc.NewGZIPDecompressor()),
	)
	registerServices(a, s)

	// start listening
	lis, err := net.Listen("tcp", a.config.Port)
	if err != nil {
		log.Fatalf("Unable to listen on %s: %v\n", a.config.Port[1:], err)
	}
	log.Println("Listening on port:", a.config.Port[1:])
	log.Fatalln(s.Serve(lis))
}

func registerServices(amp *Amplifier, s *grpc.Server) {
	var wg sync.WaitGroup
	wg.Add(len(serviceInitializers))
	for _, f := range serviceInitializers {
		go func(f serviceInitializer) {
			defer wg.Done()
			f(amp, s)
		}(f)
	}

	// Wait for all service registrations to complete.
	wg.Wait()
}

func registerVersionServer(amp *Amplifier, s *grpc.Server) {
	version.RegisterVersionServer(s, &version.Server{
		Info: &version.Info{
			Version:       amp.config.Version,
			Build:         amp.config.Build,
			GoVersion:     runtime.Version(),
			Os:            runtime.GOOS,
			Arch:          runtime.GOARCH,
			Registration:  amp.config.Registration,
			Notifications: amp.config.Notifications,
		},
	})
}

func registerLogsServer(amp *Amplifier, s *grpc.Server) {
	logs.RegisterLogsServer(s, &logs.Server{
		ES: amp.es,
		NS: amp.ns,
	})
}

func registerStatsServer(amp *Amplifier, s *grpc.Server) {
	stats.RegisterStatsServer(s, &stats.Stats{
		ES: amp.es,
	})
}

func registerAccountServer(amp *Amplifier, s *grpc.Server) {
	account.RegisterAccountServer(s, &account.Server{
		Accounts: amp.accounts,
		Config:   amp.config,
		Mailer:   amp.mailer,
		Tokens:   amp.tokens,
	})
}

func registerStackServer(amp *Amplifier, s *grpc.Server) {
	stack.RegisterStackServer(s, &stack.Server{
		Accounts: amp.accounts,
		Docker:   amp.docker,
		Stacks:   amp.stacks,
	})
}

func registerClusterServer(amp *Amplifier, s *grpc.Server) {
	cluster.RegisterClusterServer(s, &cluster.Server{
		Docker: amp.docker,
	})
}

func registerServiceServer(amp *Amplifier, s *grpc.Server) {
	service.RegisterServiceServer(s, &service.Server{
		Accounts: amp.accounts,
		Docker:   amp.docker,
		Stacks:   amp.stacks,
	})
}

func registerNodeServer(amp *Amplifier, s *grpc.Server) {
	node.RegisterNodeServer(s, &node.Server{
		Docker: amp.docker,
	})
}

func registerResourceServer(amp *Amplifier, s *grpc.Server) {
	resource.RegisterResourceServer(s, &resource.Server{
		Accounts:   amp.accounts,
		Dashboards: amp.dashboards,
		Stacks:     amp.stacks,
	})
}

func registerDashboardServer(amp *Amplifier, s *grpc.Server) {
	dashboard.RegisterDashboardServer(s, &dashboard.Server{
		Dashboards: amp.dashboards,
	})
}
