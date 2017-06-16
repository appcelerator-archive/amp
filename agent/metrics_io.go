package core

import (
	"time"

	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/docker/docker/api/types"
)

// IOStats IO stats
type IOStats struct {
	Time   time.Time
	Reads  float64
	Writes float64
	Totals float64
}

// IOStatsDiff diff between two IOStats
type IOStatsDiff struct {
	Duration float64
	Reads    float64
	Writes   float64
	Totals   float64
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
	entry.Io.Read += int64(diff.Reads)
	entry.Io.Write += int64(diff.Writes)
	entry.Io.Total += int64(diff.Totals)
}

// create new io stats
func (a *Agent) newIOStats(stats *types.StatsJSON) *IOStats {
	var io = &IOStats{Time: stats.Read}
	for _, s := range stats.BlkioStats.IoServicedRecursive {
		if s.Op == "Read" {
			io.Reads += float64(s.Value)
		} else if s.Op == "Write" {
			io.Writes += float64(s.Value)
		} else if s.Op == "Total" {
			io.Totals += float64(s.Value)
		}
	}
	return io
}

// create a new io diff computing difference between two io stats
func (a *Agent) newIODiff(newIO *IOStats, previousIO *IOStats) *IOStatsDiff {
	diff := &IOStatsDiff{Duration: newIO.Time.Sub(previousIO.Time).Minutes()}
	if diff.Duration <= 0 {
		return nil
	}
	diff.Reads = (newIO.Reads - previousIO.Reads) / diff.Duration
	diff.Writes = (newIO.Writes - previousIO.Writes) / diff.Duration
	diff.Totals = (newIO.Totals - previousIO.Totals) / diff.Duration
	return diff
}
