package core

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"time"

	"fmt"

	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/docker/docker/api/types"
	"github.com/golang/protobuf/proto"
)

// verify all containers to open metrics stream if not already done
func (a *Agent) updateMetricsStreams() {
	for ID, data := range a.containers {
		if data.metricsStream == nil || data.metricsReadError {
			streamb, err := a.dock.GetClient().ContainerStats(context.Background(), ID, true)
			if err != nil {
				log.Printf("Error opening metrics stream on container: %s\n", data.name)
			} else {
				log.Printf("open metrics stream on container: %s\n", data.name)
				data.metricsStream = streamb.Body
				go a.startReadingMetrics(ID, data)
			}
		}
	}
}

// open a metrics container stream
func (a *Agent) startReadingMetrics(ID string, data *ContainerData) {
	stream := data.metricsStream
	log.Printf("start reading metrics on container: %s\n", data.name)
	decoder := json.NewDecoder(stream)
	statsData := new(types.StatsJSON)
	var previous, now int64
	metricsEntry := &data.squashedMetricsMessage
	data.squashedMetricsMessage.ContainerId = ID
	data.squashedMetricsMessage.ContainerName = data.name
	data.squashedMetricsMessage.ContainerShortName = data.shortName
	data.squashedMetricsMessage.ContainerState = data.state
	data.squashedMetricsMessage.ServiceName = data.serviceName
	data.squashedMetricsMessage.ServiceId = data.serviceID
	data.squashedMetricsMessage.TaskId = data.taskID
	data.squashedMetricsMessage.StackName = data.stackName
	data.squashedMetricsMessage.NodeId = data.nodeID
	data.squashedMetricsMessage.Labels = data.labels
	for err := decoder.Decode(statsData); err != io.EOF; err = decoder.Decode(statsData) {
		if err != nil {
			if err.Error() == "EOF" {
				log.Printf("Stream metrics EOF container terminated: %s\n", data.name)
			} else {
				log.Printf("error reading metrics, closing metrics stream on container %s (%v)\n", data.name, err)
			}
			data.metricsReadError = true
			_ = stream.Close()
			a.removeContainer(ID)
			return
		}
		now = time.Now().UnixNano()
		if now <= previous {
			now = previous + 1
		}
		previous = now
		data.squashedMetricsMessage.Timestamp = statsData.Read.Format(time.RFC3339Nano)
		data.squashedMetricsMessage.TimeId = fmt.Sprintf("%016X", now)
		a.setMemMetrics(statsData, metricsEntry)
		a.setIOMetrics(data, statsData, metricsEntry)
		a.setNetMetrics(data, statsData, metricsEntry)
		a.setCPUMetrics(statsData, metricsEntry)
		data.squashedMetricsMessageNb++
		a.nbMetricsComputed++
	}
}

func (a *Agent) sendSquashedMetricsMessages() {
	for _, data := range a.containers {
		if data.squashedMetricsMessageNb > 0 {
			a.computeMetricsMessageAvg(data)
			a.sendMetricsMessage(data)
		}
		a.clearMetricsMessage(data)
	}
}

func (a *Agent) computeMetricsMessageAvg(data *ContainerData) {
	entry := data.squashedMetricsMessage
	entry.Cpu.TotalUsage /= float64(data.squashedMetricsMessageNb)
	entry.Cpu.UsageInKernelMode /= float64(data.squashedMetricsMessageNb)
	entry.Cpu.UsageInUserMode /= float64(data.squashedMetricsMessageNb)
	//
	entry.Io.Read /= data.squashedMetricsMessageNb
	entry.Io.Write /= data.squashedMetricsMessageNb
	entry.Io.Total /= data.squashedMetricsMessageNb
	//
	entry.Mem.Failcnt /= data.squashedMetricsMessageNb
	entry.Mem.Limit /= data.squashedMetricsMessageNb
	entry.Mem.Maxusage /= data.squashedMetricsMessageNb
	entry.Mem.Usage /= data.squashedMetricsMessageNb
	entry.Mem.UsageP /= float64(data.squashedMetricsMessageNb)
	//
	entry.Net.TotalBytes /= data.squashedMetricsMessageNb
	entry.Net.RxBytes /= data.squashedMetricsMessageNb
	entry.Net.RxDropped /= data.squashedMetricsMessageNb
	entry.Net.RxErrors /= data.squashedMetricsMessageNb
	entry.Net.RxPackets /= data.squashedMetricsMessageNb
	entry.Net.TxBytes /= data.squashedMetricsMessageNb
	entry.Net.TxDropped /= data.squashedMetricsMessageNb
	entry.Net.TxErrors /= data.squashedMetricsMessageNb
	entry.Net.TxPackets /= data.squashedMetricsMessageNb
}

func (a *Agent) sendMetricsMessage(data *ContainerData) {
	encoded, err := proto.Marshal(&data.squashedMetricsMessage)
	if err != nil {
		log.Printf("error marshalling metrics entries: %v\n", err)
		return
	}
	_, err = a.natsStreaming.GetClient().PublishAsync(ns.MetricsSubject, encoded, nil)
	if err != nil {
		log.Printf("error sending metrics entries: %v\n", err)
		return
	}
	a.nbMetrics++
}

func (a *Agent) clearMetricsMessage(data *ContainerData) {
	data.squashedMetricsMessageNb = 0
	entry := data.squashedMetricsMessage
	entry.Cpu.TotalUsage = 0
	entry.Cpu.UsageInKernelMode = 0
	entry.Cpu.UsageInUserMode = 0
	//
	entry.Io.Read = 0
	entry.Io.Write = 0
	entry.Io.Total = 0
	//
	entry.Mem.Failcnt = 0
	entry.Mem.Limit = 0
	entry.Mem.Maxusage = 0
	entry.Mem.Usage = 0
	entry.Mem.UsageP = 0
	//
	entry.Net.TotalBytes = 0
	entry.Net.RxBytes = 0
	entry.Net.RxDropped = 0
	entry.Net.RxErrors = 0
	entry.Net.RxPackets = 0
	entry.Net.TxBytes = 0
	entry.Net.TxDropped = 0
	entry.Net.TxErrors = 0
	entry.Net.TxPackets = 0
}

// close all metrics streams
func (a *Agent) closeMetricsStreams() {
	for _, data := range a.containers {
		if data.metricsStream != nil {
			err := data.metricsStream.Close()
			if err != nil {
				log.Println("Error closing a metrics stream: ", err)
			}
		}
	}
}
