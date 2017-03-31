package core

import (
	"time"

	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/docker/docker/api/types"
)

// IOStats IO stats
type IOStats struct {
	Time   time.Time
	Reads  uint64
	Writes uint64
	Totals uint64
}

// IOStatsDiff diff between two IOStats
type IOStatsDiff struct {
	Duration int64
	Reads    int64
	Writes   int64
	Totals   int64
}

// publish one IO metrics event
func (a *Agent) setIOMetrics(data *ContainerData, statsData *types.StatsJSON, entry *stats.MetricsEntry) {
	io := a.newIOStats(statsData)
	if data.previousIOStats == nil {
		data.previousIOStats = io
		return
	}
	diff := a.newIODiff(io, data.previousIOStats)
	if diff == nil {
		return
	}
	data.previousIOStats = io
	entry.Io = &stats.MetricsIOEntry{
		Read:  diff.Reads,
		Write: diff.Writes,
		Total: diff.Totals,
	}
}

// create new io stats
func (a *Agent) newIOStats(stats *types.StatsJSON) *IOStats {
	var io = &IOStats{Time: stats.Read}
	for _, s := range stats.BlkioStats.IoServicedRecursive {
		if s.Op == "Read" {
			io.Reads += s.Value
		} else if s.Op == "Write" {
			io.Writes += s.Value
		} else if s.Op == "Total" {
			io.Totals += s.Value
		}
	}
	return io
}

// create a new io diff computing difference between two io stats
func (a *Agent) newIODiff(newIO *IOStats, previousIO *IOStats) *IOStatsDiff {
	diff := &IOStatsDiff{Duration: int64(newIO.Time.Sub(previousIO.Time).Seconds())}
	if diff.Duration <= 0 {
		return nil
	}
	diff.Reads = int64(newIO.Reads - previousIO.Reads)
	diff.Writes = int64(newIO.Writes - previousIO.Writes)
	diff.Totals = int64(newIO.Totals - previousIO.Totals)
	return diff
}
