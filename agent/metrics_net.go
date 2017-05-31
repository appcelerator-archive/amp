package core

import (
	"time"

	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/docker/docker/api/types"
)

// NetStats net stats
type NetStats struct {
	Time      time.Time
	RxBytes   uint64
	RxDropped uint64
	RxErrors  uint64
	RxPackets uint64
	TxBytes   uint64
	TxDropped uint64
	TxErrors  uint64
	TxPackets uint64
}

// NetStatsDiff diff between two IOStats
type NetStatsDiff struct {
	Duration  int64
	RxBytes   int64
	RxDropped int64
	RxErrors  int64
	RxPackets int64
	TxBytes   int64
	TxDropped int64
	TxErrors  int64
	TxPackets int64
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
	entry.Net.TotalBytes += diff.RxBytes + diff.TxBytes
	entry.Net.RxBytes += diff.RxBytes
	entry.Net.RxDropped += diff.RxDropped
	entry.Net.RxErrors += diff.RxErrors
	entry.Net.RxPackets += diff.RxPackets
	entry.Net.TxBytes += diff.TxBytes
	entry.Net.TxDropped += diff.TxDropped
	entry.Net.TxErrors += diff.TxErrors
	entry.Net.TxPackets += diff.TxPackets
}

// create a new net stats
func (a *Agent) newNetStats(stats *types.StatsJSON) *NetStats {
	var net = &NetStats{Time: stats.Read}
	for _, netStats := range stats.Networks {
		net.RxBytes += netStats.RxBytes
		net.RxDropped += netStats.RxDropped
		net.RxErrors += netStats.RxErrors
		net.RxPackets += netStats.RxPackets
		net.TxBytes += netStats.TxBytes
		net.TxDropped += netStats.TxDropped
		net.TxErrors += netStats.TxErrors
		net.TxPackets += netStats.TxPackets
	}
	return net
}

// create a new net diff computing difference between two net stats
func (a *Agent) newNetDiff(newNet *NetStats, previousNet *NetStats) *NetStatsDiff {
	diff := &NetStatsDiff{Duration: int64(newNet.Time.Sub(previousNet.Time).Seconds())}
	if diff.Duration <= 0 {
		return nil
	}
	diff.RxBytes = int64(newNet.RxBytes-previousNet.RxBytes) / diff.Duration
	diff.RxDropped = int64(newNet.RxDropped-previousNet.RxDropped) / diff.Duration
	diff.RxErrors = int64(newNet.RxErrors-previousNet.RxErrors) / diff.Duration
	diff.RxPackets = int64(newNet.RxPackets-previousNet.RxPackets) / diff.Duration
	diff.TxBytes = int64(newNet.TxBytes-previousNet.TxBytes) / diff.Duration
	diff.TxDropped = int64(newNet.TxDropped-previousNet.TxDropped) / diff.Duration
	diff.TxErrors = int64(newNet.TxErrors-previousNet.TxErrors) / diff.Duration
	diff.TxPackets = int64(newNet.TxPackets-previousNet.TxPackets) / diff.Duration
	return diff
}
