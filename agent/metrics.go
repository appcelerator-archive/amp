package core

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/appcelerator/amp/api/rpc/stats"
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
	for err := decoder.Decode(statsData); err != io.EOF; err = decoder.Decode(statsData) {
		if err != nil {
			log.Printf("close metrics stream on container %s (%v)\n", data.name, err)
			data.metricsReadError = true
			_ = stream.Close()
			a.removeContainer(ID)
			return
		}
		metricsEntry := &stats.MetricsEntry{
			Timestamp:          statsData.Read.Format(time.RFC3339Nano),
			ContainerId:        ID,
			ContainerName:      data.name,
			ContainerShortName: data.shortName,
			ContainerState:     data.state,
			ServiceName:        data.serviceName,
			ServiceId:          data.serviceID,
			TaskId:             data.taskID,
			StackName:          data.stackName,
			NodeId:             data.nodeID,
			Labels:             data.labels,
		}
		a.setMemMetrics(statsData, metricsEntry)
		a.setIOMetrics(data, statsData, metricsEntry)
		a.setNetMetrics(data, statsData, metricsEntry)
		a.setCPUMetrics(statsData, metricsEntry)
		encoded, err := proto.Marshal(metricsEntry)
		if err != nil {
			log.Printf("error marshalling metrics entry: %v\n", err)
		}
		_, err = a.natsStreaming.GetClient().PublishAsync(ns.MetricsSubject, encoded, nil)
		if err != nil {
			log.Printf("error sending log entry: %v\n", err)
			return
		}
		a.nbMetrics++
	}
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
