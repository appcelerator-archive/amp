package server

import (
	"log"
	"net"
	runInfo "runtime"
	"strings"
	"time"

	// "github.com/appcelerator/amp/api/rpc/build"
	"fmt"
	"github.com/appcelerator/amp/api/rpc/function"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/oauth"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/appcelerator/amp/api/rpc/storage"
	"github.com/appcelerator/amp/api/rpc/topic"
	"github.com/appcelerator/amp/api/rpc/version"
	"github.com/appcelerator/amp/api/runtime"
	"github.com/appcelerator/amp/config"
	"github.com/appcelerator/amp/data/influx"
	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/docker/docker/client"
	"google.golang.org/grpc"
	"os"
	"sync"
)

const (
	defaultTimeOut = 30 * time.Second
)

func initDependencies(config Config) {
	// ensure all initialization code fails fast on errors; there is no point in
	// attempting to continue in a degraded state if there are problems at start up

	var wg sync.WaitGroup
	type initFunc func(Config) error

	initFuncs := []initFunc{initEtcd, initElasticsearch, initNats, initInfluxDB, initDocker}
	for _, f := range initFuncs {
		wg.Add(1)
		go func(f initFunc) {
			defer wg.Done()
			if err := f(config); err != nil {
				log.Fatalln(err)
			}
		}(f)
	}

	// Wait for all inits to complete.
	wg.Wait()
}

// Start starts the server
func Start(config Config) {
	initDependencies(config)

	// register services
	s := grpc.NewServer()
	// project.RegisterProjectServer(s, &project.Service{})
	logs.RegisterLogsServer(s, &logs.Server{
		Es:            &runtime.Elasticsearch,
		Store:         runtime.Store,
		NatsStreaming: runtime.NatsStreaming,
	})
	stats.RegisterStatsServer(s, &stats.Stats{
		Influx: runtime.Influx,
	})
	oauth.RegisterGithubServer(s, &oauth.Oauth{
		Store:        runtime.Store,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
	})
	// build.RegisterAmpBuildServer(s, &build.Proxy{})
	service.RegisterServiceServer(s, &service.Service{
		Docker: runtime.Docker,
	})
	stack.RegisterStackServiceServer(s, stack.NewServer(
		runtime.Store,
		runtime.Docker,
	))
	topic.RegisterTopicServer(s, &topic.Server{
		Store:         runtime.Store,
		NatsStreaming: runtime.NatsStreaming,
	})
	function.RegisterFunctionServer(s, &function.Server{
		Store:         runtime.Store,
		NatsStreaming: runtime.NatsStreaming,
	})
	//register storage service
	storage.RegisterStorageServer(s, &storage.Server{
		Store: runtime.Store,
	})
	version.RegisterVersionServer(s, &version.Server{
		Version:   config.Version,
		Port:      config.Port,
		GoVersion: runInfo.Version(),
		Os:        runInfo.GOOS,
		Arch:      runInfo.GOARCH,
	})

	// start listening
	lis, err := net.Listen("tcp", config.Port)
	if err != nil {
		log.Fatalf("Unable to listen on %s: %v\n", config.Port[1:], err)
	}
	log.Println("Listening on port:", config.Port[1:])
	log.Fatalln(s.Serve(lis))
}

func initEtcd(config Config) error {
	log.Println("Connecting to etcd at", strings.Join(config.EtcdEndpoints, ","))
	runtime.Store = etcd.New(config.EtcdEndpoints, "amp")
	if err := runtime.Store.Connect(defaultTimeOut); err != nil {
		return fmt.Errorf("unable to connect to etcd at %s: %v", config.EtcdEndpoints, err)
	}
	log.Println("Connected to etcd at", strings.Join(runtime.Store.Endpoints(), ","))
	return nil
}

func initElasticsearch(config Config) error {
	log.Println("Connecting to elasticsearch at", config.ElasticsearchURL)
	if err := runtime.Elasticsearch.Connect(config.ElasticsearchURL, defaultTimeOut); err != nil {
		return fmt.Errorf("unable to connect to elasticsearch at %s: %v", config.ElasticsearchURL, err)
	}
	log.Println("Connected to elasticsearch at", config.ElasticsearchURL)
	return nil
}

func initInfluxDB(config Config) error {
	log.Println("Connecting to InfluxDB at", config.InfluxURL)
	runtime.Influx = influx.New(config.InfluxURL, "telegraf", "", "")
	if err := runtime.Influx.Connect(defaultTimeOut); err != nil {
		return fmt.Errorf("unable to connect to influxDB at %s: %v", config.InfluxURL, err)
	}
	log.Println("Connected to influxDB at", config.InfluxURL)
	return nil
}

func initNats(config Config) error {
	// NATS
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("unable to get hostname: %v", err)
	}
	if runtime.NatsStreaming.Connect(config.NatsURL, amp.NatsClusterID, os.Args[0]+"-"+hostname, amp.DefaultTimeout) != nil {
		return err
	}
	return nil
}

func initDocker(config Config) error {
	log.Printf("Connecting to Docker API at %s version API: %s\n", config.DockerURL, config.DockerVersion)
	defaultHeaders := map[string]string{"User-Agent": "amplifier-1.0"}
	var err error
	runtime.Docker, err = client.NewClient(config.DockerURL, config.DockerVersion, nil, defaultHeaders)
	if err != nil {
		return fmt.Errorf("unable to connect to Docker at %s: %v", config.DockerURL, err)
	}
	log.Println("Connected to Docker API at", config.DockerURL)
	return nil
}
