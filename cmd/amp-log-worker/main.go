package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/config"
	"github.com/appcelerator/amp/data/elasticsearch"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats-streaming"
	"golang.org/x/net/context"
)

// build vars
var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string

	// es is the elasticsearch client
	es elasticsearch.Elasticsearch

	// natsStreaming is the nats streaming client
	natsStreaming ns.NatsStreaming
)

func main() {
	log.Printf("%s (version: %s, build: %s)\n", os.Args[0], Version, Build)

	err := es.Connect(amp.ElasticsearchDefaultURL, 60*time.Second)
	if err != nil {
		log.Fatalf("Unable to connect to elasticsearch on %s: %s", amp.ElasticsearchDefaultURL, err)
	}
	log.Printf("Connected to elasticsearch at %s\n", amp.ElasticsearchDefaultURL)

	es.CreateIndexIfNotExists(context.Background(), logs.EsIndex, logs.EsType, logs.EsMapping)
	if err != nil {
		log.Fatalf("Unable to create index: %s", err)
	}
	log.Printf("Created index %s\n", logs.EsIndex)

	// NATS Connect
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Unable to get hostname: %s", err)
	}
	if natsStreaming.Connect(amp.NatsDefaultURL, amp.NatsClusterID, os.Args[0]+"-"+hostname, amp.DefaultTimeout) != nil {
		log.Fatal(err)
	}

	// NATS, subscribe to function topic
	log.Println("Subscribing to topic:", amp.NatsLogsTopic)
	_, err = natsStreaming.GetClient().Subscribe(amp.NatsLogsTopic, messageHandler, stan.DeliverAllAvailable())
	if err != nil {
		natsStreaming.Close()
		log.Fatalln("Unable to subscribe to topic", err)
	}
	log.Println("Subscribed to topic:", amp.NatsLogsTopic)

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			log.Println("\nReceived an interrupt, unsubscribing and closing connection...")
			natsStreaming.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func messageHandler(msg *stan.Msg) {
	logEntry := logs.LogEntry{}
	err := proto.Unmarshal(msg.Data, &logEntry)
	if err != nil {
		log.Printf("Error unmarshalling log entry: %v", err)
	}
	timestamp, err := time.Parse(time.RFC3339Nano, logEntry.Timestamp)
	if err != nil {
		log.Printf("Error parsing timestamp: %v", err)
	}
	logEntry.Timestamp = timestamp.Format("2006-01-02T15:04:05.999")
	err = es.Index(context.Background(), logs.EsIndex, logs.EsType, logEntry) // TODO: Should we use a timeout context ?
	if err != nil {
		log.Printf("Error indexing log entry: %v", err)
	}
}
