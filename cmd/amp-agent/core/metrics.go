package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/appcelerator/amp/api/rpc/stats"
	amp "github.com/appcelerator/amp/pkg/config"
	"github.com/docker/docker/api/types"
	"github.com/gogo/protobuf/proto"
)

// verify all containers to open metrics stream if not already done
func (a *Agent) updateMetricsStreams() {
	for ID, data := range a.containers {
		if data.metricsStream == nil || data.metricsReadError {
			streamb, err := a.dockerClient.ContainerStats(context.Background(), ID, true)
			if err != nil {
				fmt.Printf("Error opening metrics stream on container: %s\n", data.name)
			} else {
				fmt.Printf("open metrics stream on container: %s\n", data.name)
				data.metricsStream = streamb.Body
				go a.startReadingMetrics(ID, data)
			}
		}
	}
}

// open a metrics container stream
func (a *Agent) startReadingMetrics(ID string, data *ContainerData) {
	stream := data.metricsStream
	fmt.Printf("start reading metrics on container: %s\n", data.name)
	decoder := json.NewDecoder(stream)
	statsData := new(types.StatsJSON)
	for err := decoder.Decode(statsData); err != io.EOF; err = decoder.Decode(statsData) {
		if err != nil {
			fmt.Printf("close metrics stream on container %s (%v)\n", data.name, err)
			data.metricsReadError = true
			stream.Close()
			a.removeContainer(ID)
			return
		}
		metricsEntry := &stats.MetricsEntry{
			Time:               statsData.Read.UnixNano(),
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
			Role:               data.role,
		}
		a.setMemMetrics(statsData, metricsEntry)
		a.setIOMetrics(data, statsData, metricsEntry)
		a.setNetMetrics(data, statsData, metricsEntry)
		a.setCPUMetrics(statsData, metricsEntry)
		encoded, err := proto.Marshal(metricsEntry)
		if err != nil {
			fmt.Printf("error marshalling metrics entry: %v", err)
		}
		_, err = a.natsStreaming.GetClient().PublishAsync(amp.NatsMetricsTopic, encoded, nil)
		if err != nil {
			fmt.Printf("error sending log entry: %v", err)
			return
		}
		a.nbMetrics++
	}
}

// close all metrics streams
func (a *Agent) closeMetricsStreams() {
	for _, data := range a.containers {
		if data.metricsStream != nil {
			data.metricsStream.Close()
		}
	}
}
