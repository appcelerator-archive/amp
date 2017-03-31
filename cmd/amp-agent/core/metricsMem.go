package core

import (
	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/docker/docker/api/types"
)

// publish one memory metrics event
func (a *Agent) setMemMetrics(statsData *types.StatsJSON, entry *stats.MetricsEntry) {
	entry.Mem = &stats.MetricsMemEntry{
		Failcnt:  int64(statsData.MemoryStats.Failcnt),
		Limit:    int64(statsData.MemoryStats.Limit),
		Maxusage: int64(statsData.MemoryStats.MaxUsage),
		Usage:    int64(statsData.MemoryStats.Usage),
		UsageP:   a.getMemUsage(statsData),
	}
}

// compute memory usage
func (a *Agent) getMemUsage(stats *types.StatsJSON) float64 {
	if stats.MemoryStats.Limit == 0 {
		return 0
	}
	return float64(stats.MemoryStats.Usage) / float64(stats.MemoryStats.Limit)
}
