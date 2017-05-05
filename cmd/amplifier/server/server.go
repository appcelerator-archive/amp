package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"sync"

	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/api/rpc/cluster"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/rpc/stats"
	kv "github.com/appcelerator/amp/api/rpc/storage"
	"github.com/appcelerator/amp/api/rpc/version"
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

type (
	clientInitializer  func(*Amplifier) error
	serviceInitializer func(*Amplifier, *grpc.Server)
)

// Amplifier represents the AMP gRPC server
type Amplifier struct {
	config *Configuration
	docker *docker.Docker
	store  storage.Interface
	es     *elasticsearch.Elasticsearch
	ns     *ns.NatsStreaming
	mailer *mail.Mailer
}

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
	registerAccountServer,
	registerClusterServer,
}

func New(config *Configuration) *Amplifier {
	amp := &Amplifier{config: config}
	initClients(amp)
	return amp
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

func initClients(amp *Amplifier) {
	// ensure all initialization code fails fast on errors; there is no point in
	// attempting to continue in a degraded state if there are problems at start up

	var wg sync.WaitGroup
	wg.Add(len(clientInitializers))
	for _, f := range clientInitializers {
		go func(f clientInitializer) {
			defer wg.Done()
			if err := f(amp); err != nil {
				log.Fatalln(err)
			}
		}(f)
	}

	// Wait for all inits to complete.
	wg.Wait()
}

func initEtcd(amp *Amplifier) error {
	amp.store = etcd.New(amp.config.EtcdEndpoints, "amp", DefaultTimeout)
	return nil
}

func initElasticsearch(amp *Amplifier) error {
	amp.es = elasticsearch.NewClient(amp.config.ElasticsearchURL, DefaultTimeout)
	return nil
}

func initNats(amp *Amplifier) error {
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("unable to get hostname: %v", err)
	}
	amp.ns = ns.NewClient(amp.config.NatsURL, ns.ClusterID, os.Args[0]+"-"+hostname, DefaultTimeout)
	return nil
}

func initDocker(amp *Amplifier) error {
	amp.docker = docker.NewClient(amp.config.DockerURL, amp.config.DockerVersion)
	log.Printf("Connecting to Docker API at %s version API: %s\n", amp.config.DockerURL, amp.config.DockerVersion)
	if err := amp.docker.Connect(); err != nil {
		return err
	}
	log.Println("Connected to Docker API at", amp.config.DockerURL)
	return nil
}

func initMailer(amp *Amplifier) error {
	amp.mailer = mail.NewMailer(amp.config.EmailKey, amp.config.EmailSender)
	return nil
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
		Es: amp.es,
		Ns: amp.ns,
	})
}

func registerStorageServer(amp *Amplifier, s *grpc.Server) {
	kv.RegisterStorageServer(s, &kv.Server{
		Store: amp.store,
	})
}

func registerStatsServer(amp *Amplifier, s *grpc.Server) {
	stats.RegisterStatsServer(s, &stats.Stats{
		Es: amp.es,
	})
}

func registerAccountServer(amp *Amplifier, s *grpc.Server) {
	account.RegisterAccountServer(s, &account.Server{
		Accounts: accounts.NewStore(amp.store),
		Mailer:   amp.mailer,
	})
}

func registerStackServer(amp *Amplifier, s *grpc.Server) {
	stack.RegisterStackServer(s, &stack.Server{
		Accounts: accounts.NewStore(amp.store),
		Docker:   amp.docker,
		Stacks:   stacks.NewStore(amp.store),
	})
}

func registerClusterServer(amp *Amplifier, s *grpc.Server) {
	cluster.RegisterClusterServer(s, &cluster.Server{
		Docker: amp.docker,
	})
}
