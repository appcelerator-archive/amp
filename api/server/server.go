package server

import (
	"log"
	"net"
	"strings"
	"time"

	// "github.com/appcelerator/amp/api/rpc/build"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/oauth"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/appcelerator/amp/api/rpc/topic"
	"github.com/appcelerator/amp/api/runtime"
	"github.com/appcelerator/amp/config"
	"github.com/appcelerator/amp/data/influx"
	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/docker/docker/client"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/nats"
	"google.golang.org/grpc"
)

const (
	defaultTimeOut = 30 * time.Second
	natsClientID   = "amplifier"
)

func initDependencies(config Config) error {
	// ensure all initialization code fails fast on errors; there is no point in
	// attempting to continue in a degraded state if there are problems at start up
	if err := initEtcd(config); err != nil {
		return err
	}
	if err := initElasticsearch(config); err != nil {
		return err
	}
	if err := initNats(config); err != nil {
		return err
	}
	if err := initInfluxDB(config); err != nil {
		return err
	}
	if err := initDocker(config); err != nil {
		return err
	}
	return nil
}

// Start starts the server
func Start(config Config) {
	if err := initDependencies(config); err != nil {
		panic(err)
	}

	// register services
	s := grpc.NewServer()
	// project.RegisterProjectServer(s, &project.Service{})
	logs.RegisterLogsServer(s, &logs.Server{
		Es:    &runtime.Elasticsearch,
		Store: runtime.Store,
		Nats:  runtime.Nats,
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
		Store: runtime.Store,
		Nats:  runtime.Nats,
	})

	// start listening
	lis, err := net.Listen("tcp", config.Port)
	if err != nil {
		log.Fatalf("amplifer is unable to listen on: %s\n%v", config.Port[1:], err)
	}
	log.Printf("amplifier is listening on port %s\n", config.Port[1:])
	s.Serve(lis)
}

func initEtcd(config Config) error {
	log.Printf("connecting to etcd at %v", strings.Join(config.EtcdEndpoints, ","))
	runtime.Store = etcd.New(config.EtcdEndpoints, "amp")
	if err := runtime.Store.Connect(defaultTimeOut); err != nil {
		return fmt.Errorf("amplifer is unable to connect to etcd on: %s\n%v", config.EtcdEndpoints, err)
	}
	log.Printf("connected to etcd at %v", strings.Join(runtime.Store.Endpoints(), ","))
	return nil
}

func initElasticsearch(config Config) error {
	log.Printf("connecting to elasticsearch at %s\n", config.ElasticsearchURL)
	if err := runtime.Elasticsearch.Connect(config.ElasticsearchURL, defaultTimeOut); err != nil {
		return fmt.Errorf("amplifer is unable to connect to elasticsearch on: %s\n%v", config.ElasticsearchURL, err)
	}
	log.Printf("connected to elasticsearch at %s\n", config.ElasticsearchURL)
	return nil
}

func initInfluxDB(config Config) error {
	log.Printf("connecting to InfluxDB at %s\n", config.InfluxURL)
	runtime.Influx = influx.New(config.InfluxURL, "telegraf", "", "")
	if err := runtime.Influx.Connect(defaultTimeOut); err != nil {
		return fmt.Errorf("amplifer is unable to connect to influxDB on: %s\n%v", config.InfluxURL, err)
	}
	log.Printf("connected to influxDB at %s\n", config.InfluxURL)
	return nil
}

func initNats(config Config) error {
	log.Printf("Connecting to NATS-Streaming at %s\n", config.NatsURL)
	var err error

	nc, err := nats.Connect(config.NatsURL, nats.Timeout(defaultTimeOut))
	if err != nil {
		fmt.Errorf("amplifer is unable to connect to NATS on: %s\n%v", config.NatsURL, err)
	}

	runtime.Nats, err = stan.Connect(amp.NatsClusterID, natsClientID+strconv.Itoa(rand.Int()), stan.NatsConn(nc), stan.ConnectWait(defaultTimeOut))
	if err != nil {
		return fmt.Errorf("amplifer is unable to connect to NATS-Streaming on: %s\n%v", config.NatsURL, err)
	}
	log.Printf("Connected to NATS-Streaming at %s\n", config.NatsURL)
	return nil
}

func initDocker(config Config) error {
	log.Printf("connecting to Docker API at %s version API: %s\n", config.DockerURL, config.DockerVersion)
	defaultHeaders := map[string]string{"User-Agent": "amplifier-1.0"}
	var err error
	runtime.Docker, err = client.NewClient(config.DockerURL, config.DockerVersion, nil, defaultHeaders)
	if err != nil {
		return fmt.Errorf("amplifer is unable to connect to Docker on: %s\n%v", config.DockerURL, err)
	}
	log.Printf("connected to Docker at %s\n", config.DockerURL)
	return nil
}
