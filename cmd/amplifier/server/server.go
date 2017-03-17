package server

import (
	"fmt"
	"log"
	"net"
	"os"
	rt "runtime"
	"strings"
	"sync"
	"time"

	"github.com/appcelerator/amp/api/rpc/account"
// TODO: @bquenin
//	"github.com/appcelerator/amp/api/rpc/function"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/oauth"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/appcelerator/amp/api/rpc/storage"
	"github.com/appcelerator/amp/api/rpc/topic"
	"github.com/appcelerator/amp/api/rpc/version"
	"github.com/appcelerator/amp/api/runtime"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/influx"
	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/appcelerator/amp/pkg/config" // needed for amp.Config (TODO: fix this!)
	"github.com/docker/docker/client"
	"google.golang.org/grpc"
)

const (
	defaultTimeOut = 30 * time.Second
)

type (
	clientInitializer  func(*amp.Config) error
	serviceInitializer func(*amp.Config, *grpc.Server)
)

// Client initializers open connections to required backend services
// Clients are stored as members of runtime
var clientInitializers = []clientInitializer{
	//initEtcd,
	//initElasticsearch,
	//initNats,
	//initInfluxDB,
	//initDocker,
}

// Service initializers register the services with the grpc server
var serviceInitializers = []serviceInitializer{
	registerVersionServer,
	registerStorageServer,
	//registerLogsServer,
	//registerStatsServer,
	//registerServiceServer,
	//registerStackServiceServer,
	//registerTopicServer,
	//registerFunctionServer,
	//registerGithubServer,
	registerAccountServer,
}

// Start starts the amplifier server
func Start(c *amp.Config) {
	// initialize clients
	initClients(c)

	// register services
	s := grpc.NewServer()
	registerServices(c, s)

	// start listening
	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalf("Unable to listen on %s: %v\n", c.Port[1:], err)
	}
	log.Println("Listening on port:", c.Port[1:])
	log.Fatalln(s.Serve(lis))
}

func initClients(c *amp.Config) {
	// ensure all initialization code fails fast on errors; there is no point in
	// attempting to continue in a degraded state if there are problems at start up

	var wg sync.WaitGroup
	wg.Add(len(clientInitializers))
	for _, f := range clientInitializers {
		go func(f clientInitializer) {
			defer wg.Done()
			if err := f(c); err != nil {
				log.Fatalln(err)
			}
		}(f)
	}

	// Wait for all inits to complete.
	wg.Wait()
}

func initEtcd(c *amp.Config) error {
	log.Println("Connecting to etcd at", strings.Join(c.EtcdEndpoints, ","))
	runtime.Store = etcd.New(c.EtcdEndpoints, "amp")
	if err := runtime.Store.Connect(defaultTimeOut); err != nil {
		return fmt.Errorf("unable to connect to etcd at %s: %v", c.EtcdEndpoints, err)
	}
	log.Println("Connected to etcd at", strings.Join(runtime.Store.Endpoints(), ","))
	return nil
}

func initElasticsearch(c *amp.Config) error {
	log.Println("Connecting to elasticsearch at", c.ElasticsearchURL)
	if err := runtime.Elasticsearch.Connect(c.ElasticsearchURL, defaultTimeOut); err != nil {
		return fmt.Errorf("unable to connect to elasticsearch at %s: %v", c.ElasticsearchURL, err)
	}
	log.Println("Connected to elasticsearch at", c.ElasticsearchURL)
	return nil
}

func initInfluxDB(c *amp.Config) error {
	log.Println("Connecting to InfluxDB at", c.InfluxURL)
	runtime.Influx = influx.New(c.InfluxURL, "telegraf", "", "")
	if err := runtime.Influx.Connect(defaultTimeOut); err != nil {
		return fmt.Errorf("unable to connect to influxDB at %s: %v", c.InfluxURL, err)
	}
	log.Println("Connected to influxDB at", c.InfluxURL)
	return nil
}

func initNats(c *amp.Config) error {
	// NATS
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("unable to get hostname: %v", err)
	}
	if runtime.NatsStreaming.Connect(c.NatsURL, amp.NatsClusterID, os.Args[0]+"-"+hostname, amp.DefaultTimeout) != nil {
		return err
	}
	return nil
}

func initDocker(c *amp.Config) error {
	log.Printf("Connecting to Docker API at %s version API: %s\n", c.DockerURL, c.DockerVersion)
	defaultHeaders := map[string]string{"User-Agent": "amplifier-1.0"}
	var err error
	runtime.Docker, err = client.NewClient(c.DockerURL, c.DockerVersion, nil, defaultHeaders)
	if err != nil {
		return fmt.Errorf("unable to connect to Docker at %s: %v", c.DockerURL, err)
	}
	log.Println("Connected to Docker API at", c.DockerURL)
	return nil
}

func registerServices(c *amp.Config, s *grpc.Server) {
	var wg sync.WaitGroup
	for _, f := range serviceInitializers {
		wg.Add(1)
		go func(f serviceInitializer) {
			defer wg.Done()
			f(c, s)
		}(f)
	}

	// Wait for all service registrations to complete.
	wg.Wait()
}

func registerVersionServer(c *amp.Config, s *grpc.Server) {
	version.RegisterVersionServer(s, &version.Server{
		Version:   c.Version,
		Port:      c.Port,
		GoVersion: rt.Version(),
		Os:        rt.GOOS,
		Arch:      rt.GOARCH,
	})
}

func registerLogsServer(c *amp.Config, s *grpc.Server) {
	logs.RegisterLogsServer(s, &logs.Server{
		Es:            &runtime.Elasticsearch,
		Store:         runtime.Store,
		NatsStreaming: runtime.NatsStreaming,
	})
}

func registerStorageServer(c *amp.Config, s *grpc.Server) {
	storage.RegisterStorageServer(s, &storage.Server{
		Store: runtime.Store,
	})
}

func registerStatsServer(c *amp.Config, s *grpc.Server) {
	stats.RegisterStatsServer(s, &stats.Stats{
		Influx: runtime.Influx,
	})
}

func registerServiceServer(c *amp.Config, s *grpc.Server) {
	service.RegisterServiceServer(s, &service.Service{
		Docker: runtime.Docker,
	})
}

func registerStackServiceServer(c *amp.Config, s *grpc.Server) {
	stack.RegisterStackServiceServer(s, stack.NewServer(
		runtime.Store,
		runtime.Docker,
	))
}

func registerTopicServer(c *amp.Config, s *grpc.Server) {
	topic.RegisterTopicServer(s, &topic.Server{
		Store:         runtime.Store,
		NatsStreaming: runtime.NatsStreaming,
	})
}

func registerFunctionServer(c *amp.Config, s *grpc.Server) {
	// TODO: @bquenin needs to update
	//function.RegisterFunctionServer(s, &function.Server{
	//	Store:         runtime.Store,
	//	NatsStreaming: runtime.NatsStreaming,
	//})
}

func registerGithubServer(c *amp.Config, s *grpc.Server) {
	oauth.RegisterGithubServer(s, &oauth.Oauth{
		Store:        runtime.Store,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
	})
}

func registerAccountServer(c *amp.Config, s *grpc.Server) {
	account.RegisterAccountServer(s, &account.Server{
		Accounts: accounts.NewStore(runtime.Store),
	})
}
