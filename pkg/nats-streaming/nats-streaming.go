package ns

import (
	"fmt"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/nats"
	"log"
	"time"
)

// NatsStreaming NATS-Streaming wrapper
type NatsStreaming struct {
	client stan.Conn
}

// Connect to NATS-Streaming
func (ns *NatsStreaming) Connect(url string, clusterID string, clientID string, timeout time.Duration) error {
	log.Printf("Connecting to nats streaming, url: %s, clusterId: %s, clientId: %s, timeout: %s\n", url, clusterID, clientID, timeout)
	nc, err := nats.Connect(url, nats.Timeout(timeout))
	if err != nil {
		return fmt.Errorf("unable to connect to nats streaming: %v", err)
	}

	ns.client, err = stan.Connect(clusterID, clientID, stan.NatsConn(nc), stan.ConnectWait(timeout))
	if err != nil {
		return fmt.Errorf("unable to connect to nats streaming: %v", err)
	}
	log.Println("Connected to nats streaming successfuly")
	return nil
}

// GetClient return client
func (ns *NatsStreaming) GetClient() stan.Conn {
	return ns.client
}

// Close the client
func (ns *NatsStreaming) Close() error {
	return ns.client.Close()
}
