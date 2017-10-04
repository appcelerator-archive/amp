package core

import (
	"docker.io/go-docker/api/types"
	"github.com/appcelerator/amp/api/rpc/stats"
)

// publish one memory metrics event
func (a *Agent) setMemMetrics(statsData *types.StatsJSON, entry *stats.MetricsEntry) {
	entry.Mem.Failcnt += int64(statsData.MemoryStats.Failcnt)
	entry.Mem.Limit += int64(statsData.MemoryStats.Limit)
	entry.Mem.Maxusage += int64(statsData.MemoryStats.MaxUsage)
	entry.Mem.Usage += int64(statsData.MemoryStats.Usage)
	entry.Mem.UsageP += a.getMemUsage(statsData)
}

// compute memory usage
func (a *Agent) getMemUsage(stats *types.StatsJSON) float64 {
	if stats.MemoryStats.Limit == 0 {
		return 0
	}
	return float64(stats.MemoryStats.Usage) / float64(stats.MemoryStats.Limit)
}
