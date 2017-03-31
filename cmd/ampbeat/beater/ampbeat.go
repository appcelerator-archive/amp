package beater

import (
	"fmt"
	"os"
	"time"

	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/appcelerator/amp/cmd/ampbeat/config"
	"github.com/appcelerator/amp/pkg/config"
	ns "github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats-streaming"
)

type Ampbeat struct {
	done          chan struct{}
	config        config.Config
	client        publisher.Client
	natsStreaming ns.NatsStreaming
}

var bt = &Ampbeat{
	done: make(chan struct{}),
}

// Creates beater
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
	if err = bt.natsStreaming.Connect(amp.NatsDefaultURL, amp.NatsClusterID, b.Name+"-"+hostname, amp.DefaultTimeout); err != nil {
		return nil, fmt.Errorf("Unable to connect to NATS: %v", err)
	}

	return bt, nil
}

func (bt *Ampbeat) Run(b *beat.Beat) error {
	logp.Info("ampbeat is running! Hit CTRL-C to stop it.")
	bt.client = b.Publisher.Connect()

	// logs subscription
	if _, err := bt.natsStreaming.GetClient().Subscribe(amp.NatsLogsTopic, logMessageHandler, stan.DeliverAllAvailable()); err != nil {
		return fmt.Errorf("Unable to subscribe to subject: %v", err)
	}
	logp.Info("Succesfully subscribed to logs subject")

	// metrics subscription
	if _, err := bt.natsStreaming.GetClient().Subscribe(amp.NatsMetricsTopic, metricsMessageHandler, stan.DeliverAllAvailable()); err != nil {
		return fmt.Errorf("Unable to subscribe to subject: %v", err)
	}
	logp.Info("Succesfully subscribed to metrics subject")

	select {
	case <-bt.done:
		return nil
	}
	return nil
}

func (bt *Ampbeat) Stop() {
	bt.natsStreaming.Close()
	bt.client.Close()
	close(bt.done)
}

func logMessageHandler(msg *stan.Msg) {
	e := logs.LogEntry{}
	err := proto.Unmarshal(msg.Data, &e)
	if err != nil {
		logp.Err("Error unmarshalling log entry: %s", err.Error())
		return
	}
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
		"stack_name":           e.StackName,
		"node_id":              e.NodeId,
		"role":                 e.Role,
		"message":              e.Message,
	}
	bt.client.PublishEvent(event)
}

func metricsMessageHandler(msg *stan.Msg) {
	e := stats.MetricsEntry{}
	err := proto.Unmarshal(msg.Data, &e)
	if err != nil {
		logp.Err("Error unmarshalling metrics entry: %s", err.Error())
		return
	}
	timestamp, err := time.Parse(time.RFC3339Nano, e.Timestamp)
	if err != nil {
		timestamp = time.Now()
	}
	event := common.MapStr{
		"@timestamp":           common.Time(timestamp),
		"type":                 "metrics",
		"container_id":         e.ContainerId,
		"container_name":       e.ContainerName,
		"container_short_name": e.ContainerShortName,
		"container_state":      e.ContainerState,
		"service_name":         e.ServiceName,
		"service_id":           e.ServiceId,
		"task_id":              e.TaskId,
		"stack_name":           e.StackName,
		"node_id":              e.NodeId,
		"role":                 e.Role,
	}
	if e.Cpu != nil {
		event["cpu"] = common.MapStr{
			"totalUsage":        e.Cpu.TotalUsage,
			"usageInKernelmode": e.Cpu.UsageInKernelMode,
			"usageInUsermode":   e.Cpu.UsageInUserMode,
		}
	}
	if e.Io != nil {
		event["io"] = common.MapStr{
			"read":  e.Io.Read,
			"write": e.Io.Write,
			"total": e.Io.Total,
		}
	}
	if e.Mem != nil {
		event["mem"] = common.MapStr{
			"failcnt":  e.Mem.Failcnt,
			"limit":    e.Mem.Limit,
			"maxUsage": e.Mem.Maxusage,
			"usage":    e.Mem.Usage,
			"usage_p":  e.Mem.UsageP,
		}
	}
	if e.Net != nil {
		event["net"] = common.MapStr{
			"totalBytes": e.Net.TotalBytes,
			"rxBytes":    e.Net.RxBytes,
			"rxDropped":  e.Net.RxDropped,
			"rxErrors":   e.Net.RxErrors,
			"rxPackets":  e.Net.RxPackets,
			"txBytes":    e.Net.TxBytes,
			"txDropped":  e.Net.TxDropped,
			"txErrors":   e.Net.TxErrors,
			"txPackets":  e.Net.TxPackets,
		}
	}
	bt.client.PublishEvent(event)
}
