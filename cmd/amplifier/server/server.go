package server

import (
	"log"
	"net"
	"os"
	"runtime"
	"sync"

	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/api/rpc/cluster"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/node"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/rpc/stats"
	kv "github.com/appcelerator/amp/api/rpc/storage"
	"github.com/appcelerator/amp/api/rpc/version"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/appcelerator/amp/data/accounts"
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
	config   *configuration.Configuration
	docker   *docker.Docker
	store    storage.Interface
	es       *elasticsearch.Elasticsearch
	ns       *ns.NatsStreaming
	mailer   *mail.Mailer
	accounts accounts.Interface
}

// Service initializers register the services with the grpc server
var serviceInitializers = []serviceInitializer{
	registerVersionServer,
	registerStorageServer,
	registerLogsServer,
	registerStackServer,
	registerStatsServer,
	registerAccountServer,
	registerClusterServer,
	registerServiceServer,
	registerNodeServer,
}

func New(config *configuration.Configuration) (*Amplifier, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	etcd := etcd.New(config.EtcdEndpoints, "amp", configuration.DefaultTimeout)
	amp := &Amplifier{
		config:   config,
		store:    etcd,
		es:       elasticsearch.NewClient(config.ElasticsearchURL, configuration.DefaultTimeout),
		ns:       ns.NewClient(config.NatsURL, ns.ClusterID, "amplifier-"+hostname, configuration.DefaultTimeout),
		docker:   docker.NewClient(config.DockerURL, config.DockerVersion),
		mailer:   mail.NewMailer(config.EmailKey, config.EmailSender, config.Notifications),
		accounts: accounts.NewStore(etcd, config.Registration),
	}
	return amp, nil
}

// Start starts the amplifier server
func (a *Amplifier) Start() {
	// register services
	s := grpc.NewServer(
		grpc.StreamInterceptor(auth.StreamInterceptor),
		grpc.UnaryInterceptor(auth.Interceptor),
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
			Version:   amp.config.Version,
			Build:     amp.config.Build,
			GoVersion: runtime.Version(),
			Os:        runtime.GOOS,
			Arch:      runtime.GOARCH,
		},
	})
}

func registerLogsServer(amp *Amplifier, s *grpc.Server) {
	logs.RegisterLogsServer(s, &logs.Server{
		ES: amp.es,
		NS: amp.ns,
	})
}

func registerStorageServer(amp *Amplifier, s *grpc.Server) {
	kv.RegisterStorageServer(s, &kv.Server{
		Store: amp.store,
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
		Mailer:   amp.mailer,
		Config:   amp.config,
	})
}

func registerStackServer(amp *Amplifier, s *grpc.Server) {
	stack.RegisterStackServer(s, &stack.Server{
		Accounts: amp.accounts,
		Docker:   amp.docker,
		Stacks:   stacks.NewStore(amp.store, amp.accounts),
	})
}

func registerClusterServer(amp *Amplifier, s *grpc.Server) {
	cluster.RegisterClusterServer(s, &cluster.Server{
		Docker: amp.docker,
	})
}

func registerServiceServer(amp *Amplifier, s *grpc.Server) {
	service.RegisterServiceServer(s, &service.Server{
		Docker: amp.docker,
	})
}

func registerServiceServer(c *Configuration, s *grpc.Server) {
	service.RegisterServiceServer(s, &service.Server{
		Docker: runtime.Docker,
	})
}

func registerNodeServer(c *Configuration, s *grpc.Server) {
	node.RegisterNodeServer(s, &node.Server{
		Docker: runtime.Docker,
	})
}
