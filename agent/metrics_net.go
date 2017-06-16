package core

import (
	"time"

	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/docker/docker/api/types"
)

// NetStats net stats
type NetStats struct {
	Time      time.Time
	RxBytes   float64
	RxDropped float64
	RxErrors  float64
	RxPackets float64
	TxBytes   float64
	TxDropped float64
	TxErrors  float64
	TxPackets float64
}

// NetStatsDiff diff between two IOStats
type NetStatsDiff struct {
	Duration  float64
	RxBytes   float64
	RxDropped float64
	RxErrors  float64
	RxPackets float64
	TxBytes   float64
	TxDropped float64
	TxErrors  float64
	TxPackets float64
}

// publish one net metrics event
func (a *Agent) setNetMetrics(data *ContainerData, statsData *types.StatsJSON, entry *stats.MetricsEntry) {
	net := a.newNetStats(statsData)
	if data.previousNetStats == nil {
		data.previousNetStats = net
		return
	}
	diff := a.newNetDiff(net, data.previousNetStats)
	if diff == nil {
		return
	}
	data.previousNetStats = net
	entry.Net.TotalBytes += int64(diff.RxBytes + diff.TxBytes)
	entry.Net.RxBytes += int64(diff.RxBytes)
	entry.Net.RxDropped += int64(diff.RxDropped)
	entry.Net.RxErrors += int64(diff.RxErrors)
	entry.Net.RxPackets += int64(diff.RxPackets)
	entry.Net.TxBytes += int64(diff.TxBytes)
	entry.Net.TxDropped += int64(diff.TxDropped)
	entry.Net.TxErrors += int64(diff.TxErrors)
	entry.Net.TxPackets += int64(diff.TxPackets)
}

// create a new net stats
func (a *Agent) newNetStats(stats *types.StatsJSON) *NetStats {
	var net = &NetStats{Time: stats.Read}
	for _, netStats := range stats.Networks {
		net.RxBytes += float64(netStats.RxBytes)
		net.RxDropped += float64(netStats.RxDropped)
		net.RxErrors += float64(netStats.RxErrors)
		net.RxPackets += float64(netStats.RxPackets)
		net.TxBytes += float64(netStats.TxBytes)
		net.TxDropped += float64(netStats.TxDropped)
		net.TxErrors += float64(netStats.TxErrors)
		net.TxPackets += float64(netStats.TxPackets)
	}
	return net
}

// create a new net diff computing difference between two net stats
func (a *Agent) newNetDiff(newNet *NetStats, previousNet *NetStats) *NetStatsDiff {
	diff := &NetStatsDiff{Duration: newNet.Time.Sub(previousNet.Time).Minutes()}
	if diff.Duration <= 0 {
		return nil
	}
	diff.RxBytes = (newNet.RxBytes - previousNet.RxBytes) / diff.Duration
	diff.RxDropped = (newNet.RxDropped - previousNet.RxDropped) / diff.Duration
	diff.RxErrors = (newNet.RxErrors - previousNet.RxErrors) / diff.Duration
	diff.RxPackets = (newNet.RxPackets - previousNet.RxPackets) / diff.Duration
	diff.TxBytes = (newNet.TxBytes - previousNet.TxBytes) / diff.Duration
	diff.TxDropped = (newNet.TxDropped - previousNet.TxDropped) / diff.Duration
	diff.TxErrors = (newNet.TxErrors - previousNet.TxErrors) / diff.Duration
	diff.TxPackets = (newNet.TxPackets - previousNet.TxPackets) / diff.Duration
	return diff
}
