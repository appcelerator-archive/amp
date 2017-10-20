package beater

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/cmd/ampbeat/config"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats-streaming"
)

// Ampbeat amp libbeat type
type Ampbeat struct {
	done          chan struct{}
	config        config.Config
	client        publisher.Client
	natsStreaming *ns.NatsStreaming
}

var bt = &Ampbeat{
	done: make(chan struct{}),
}

// New Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	bt.config = config.DefaultConfig
	if err := cfg.Unpack(&bt.config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	// NATS Connect
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("Unable to get hostname: %v", err)
	}
	bt.natsStreaming = ns.NewClient(ns.DefaultURL, ns.ClusterID, b.Name+"-"+hostname, time.Minute)
	if err = bt.natsStreaming.Connect(); err != nil {
		return nil, fmt.Errorf("Unable to connect to NATS: %v", err)
	}

	return bt, nil
}

//Run run the beat
func (bt *Ampbeat) Run(b *beat.Beat) error {
	logp.Info("ampbeat is running! Hit CTRL-C to stop it.")
	bt.client = b.Publisher.Connect()

	// logs subscription
	if _, err := bt.natsStreaming.GetClient().QueueSubscribe(ns.LogsSubject, ns.LogsQGroup, logMessageHandler, stan.DurableName("logs")); err != nil {
		return fmt.Errorf("Unable to subscribe to subject: %v", err)
	}
	logp.Info("Successfully subscribed to logs subject")

	select {
	case <-bt.done:
		return nil
	}
}

//Stop stop the beat
func (bt *Ampbeat) Stop() {
	bt.natsStreaming.Close()
	bt.client.Close()
	close(bt.done)
}

func logMessageHandler(msg *stan.Msg) {
	list := logs.GetReply{}
	err := proto.Unmarshal(msg.Data, &list)
	if err != nil {
		logp.Err("Error unmarshalling log entry: %s", err.Error())
		return
	}
	for _, e := range list.Entries {
		timestamp, err := time.Parse(time.RFC3339Nano, e.Timestamp)
		if err != nil {
			timestamp = time.Now()
		}
		event := common.MapStr{
			"@timestamp":           common.Time(timestamp),
			"type":                 "logs",
			"container_id":         e.ContainerId,
			"container_name":       e.ContainerName,
			"container_short_name": e.ContainerShortName,
			"container_state":      e.ContainerState,
			"service_name":         e.ServiceName,
			"service_id":           e.ServiceId,
			"task_id":              e.TaskId,
			"task_slot":            e.TaskSlot,
			"stack_name":           e.StackName,
			"stack_id":             e.StackId,
			"node_id":              e.NodeId,
			"time_id":              e.TimeId,
			"msg":                  e.Msg,
		}
		if len(e.Labels) > 0 {
			event["labels"] = getLabels(e.Labels)
		}
		bt.client.PublishEvent(event)
	}
}

// Need to replace '.'' by '-'', because '.' is used by ES for object operator it can't be part of a field name.
func getLabels(labels map[string]string) common.MapStr {
	labmap := common.MapStr{}
	for name, value := range labels {
		if !strings.HasPrefix(name, "com.docker.") {
			name = strings.Replace(name, ".", "-", -1)
			labmap[name] = value
		}
	}
	return labmap
}
