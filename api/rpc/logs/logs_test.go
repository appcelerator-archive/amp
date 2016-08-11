package logs

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/appcelerator/amp/api/server"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	defaultPort             = ":50101"
	etcdDefaultEndpoints    = "http://localhost:2379"
	serverAddress           = "localhost" + defaultPort
	elasticsearchDefaultURL = "http://localhost:9200"
)

var (
	config           server.Config
	port             string
	etcdEndpoints    string
	elasticsearchURL string
)

func parseEnv() {
	port = os.Getenv("port")
	if port == "" {
		port = defaultPort
	}
	etcdEndpoints = os.Getenv("endpoints")
	if etcdEndpoints == "" {
		etcdEndpoints = etcdDefaultEndpoints
	}
	elasticsearchURL = os.Getenv("elasticsearchURL")
	if elasticsearchURL == "" {
		elasticsearchURL = elasticsearchDefaultURL
	}

	// update config
	config.Port = port
	for _, s := range strings.Split(etcdEndpoints, ",") {
		config.EtcdEndpoints = append(config.EtcdEndpoints, s)
	}
	config.ElasticsearchURL = elasticsearchURL
}

func TestMain(m *testing.M) {
	parseEnv()
	go server.Start(config)

	// there is no event when the server starts listening, so we just wait a second
	time.Sleep(1 * time.Second)

	os.Exit(m.Run())
}

func TestShouldGetAHundredLogEntries(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}

	// Contact the server and print out its response.
	c := NewLogsClient(conn)
	r, err := c.Get(context.Background(), &GetRequest{})
	if err != nil {
		t.Fatalf("could not get logs: %v", err)
	}
	if len(r.Entries) != 100 {
		t.Fatalf("expected: %v, got: %v", 100, len(r.Entries))

	}
	conn.Close()
}
