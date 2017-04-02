package core

import (
	"time"

	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/docker/docker/api/types"
)

// CPUStats cpu stats
type CPUStats struct {
	Time                 time.Time
	PerCPUUsage          []uint64
	TotalUsage           uint64
	UsageInKernelmode    uint64
	UsageInUsermode      uint64
	PreTime              time.Time
	PrePerCPUUsage       []uint64
	PreTotalUsage        uint64
	PreUsageInKernelmode uint64
	PreUsageInUsermode   uint64
}

// CPUStatsDiff diff between two cpu stats
type CPUStatsDiff struct {
	Duration          uint64
	TotalUsage        float64
	UsageInKernelmode float64
	UsageInUsermode   float64
}

// publish one cpu event
func (a *Agent) setCPUMetrics(statsData *types.StatsJSON, entry *stats.MetricsEntry) {
	cpu := a.newCPUStats(statsData)
	diff := a.newCPUDiff(cpu)
	if diff == nil {
		return
	}
	entry.Cpu = &stats.MetricsCPUEntry{
		TotalUsage:        diff.TotalUsage,
		UsageInKernelMode: diff.UsageInKernelmode,
		UsageInUserMode:   diff.UsageInUsermode,
	}
}

// build a new cpu metrics stats
func (a *Agent) newCPUStats(stats *types.StatsJSON) *CPUStats {
	var cpu = &CPUStats{
		Time:                 stats.Read,
		PerCPUUsage:          stats.CPUStats.CPUUsage.PercpuUsage,
		TotalUsage:           stats.CPUStats.CPUUsage.TotalUsage,
		UsageInKernelmode:    stats.CPUStats.CPUUsage.UsageInKernelmode,
		UsageInUsermode:      stats.CPUStats.CPUUsage.UsageInUsermode,
		PreTime:              stats.PreRead,
		PrePerCPUUsage:       stats.PreCPUStats.CPUUsage.PercpuUsage,
		PreTotalUsage:        stats.PreCPUStats.CPUUsage.TotalUsage,
		PreUsageInKernelmode: stats.PreCPUStats.CPUUsage.UsageInKernelmode,
		PreUsageInUsermode:   stats.PreCPUStats.CPUUsage.UsageInUsermode,
	}
	return cpu
}

// build a new diff computing difference between two cpu stats
func (a *Agent) newCPUDiff(cpu *CPUStats) *CPUStatsDiff {
	diff := &CPUStatsDiff{Duration: uint64(cpu.Time.Sub(cpu.PreTime).Seconds())}
	if diff.Duration <= 0 {
		return nil
	}
	diff.TotalUsage = a.calculateLoad(cpu.TotalUsage, cpu.PreTotalUsage, diff.Duration)
	diff.UsageInKernelmode = a.calculateLoad(cpu.UsageInKernelmode, cpu.UsageInKernelmode, diff.Duration)
	diff.UsageInUsermode = a.calculateLoad(cpu.UsageInUsermode, cpu.UsageInUsermode, diff.Duration)
	return diff
}

// compute cpu usage concidering event duration
func (a *Agent) calculateLoad(oldValue uint64, newValue uint64, duration uint64) float64 {
	value := int64(oldValue - newValue)
	if value < 0 || duration == 0 {
		return float64(0)
	}
	return float64(value) / (float64(duration) * float64(10000000))
}
