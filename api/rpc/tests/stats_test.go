package tests

import (
	"strings"
	"testing"

	"github.com/appcelerator/amp/api/rpc/stats"
)

const (
	queryVal = "measurements"
)

func TestStatsQueryService(t *testing.T) {
	query := stats.StatsRequest{}
	query.Discriminator = "service"
	query.StatsCpu = true
	query.StatsMem = true
	query.StatsIo = true
	query.StatsNet = true
	query.Period = "5m"
	query.FilterServiceName = "amp"
	res, err := statsClient.StatsQuery(ctx, &query)
	if err != nil {
		t.Error(err)
	}
	if len(res.Entries) == 0 {
		t.Errorf("Unexpected empty answer from server")
	}
	for _, result := range res.Entries {
		if !strings.HasPrefix(result.ServiceName, "amp") {
			t.Errorf("Unexpected service selected: %s\n", result.ServiceName)
		}
	}
}

func TestStatsQueryContainer(t *testing.T) {
	query := stats.StatsRequest{}
	query.Discriminator = "container"
	query.StatsCpu = true
	query.StatsMem = true
	query.StatsIo = true
	query.StatsNet = true
	query.Period = "5m"
	query.FilterContainerName = "amp"
	res, err := statsClient.StatsQuery(ctx, &query)
	if err != nil {
		t.Error(err)
	}
	if len(res.Entries) == 0 {
		t.Errorf("Unexpected empty answer from server")
	}
	for _, result := range res.Entries {
		if !strings.HasPrefix(result.ContainerName, "amp") {
			t.Errorf("Unexpected container selected: %s\n", result.ContainerName)
		}
	}
}

func TestStatsQueryTask(t *testing.T) {
	query := stats.StatsRequest{}
	query.Discriminator = "task"
	query.StatsCpu = true
	query.StatsMem = true
	query.StatsIo = true
	query.StatsNet = true
	query.Period = "5m"
	query.FilterTaskName = "amp"
	res, err := statsClient.StatsQuery(ctx, &query)
	if err != nil {
		t.Error(err)
	}
	if len(res.Entries) == 0 {
		t.Errorf("Unexpected empty answer from server")
	}
	for _, result := range res.Entries {
		if !strings.HasPrefix(result.TaskName, "amp") {
			t.Errorf("Unexpected task selected: %s\n", result.TaskName)
		}
	}
}

func TestStatsQueryNode(t *testing.T) {
	query := stats.StatsRequest{}
	query.Discriminator = "node"
	query.StatsCpu = true
	query.StatsMem = true
	query.StatsIo = true
	query.StatsNet = true
	query.Period = "5m"
	res, err := statsClient.StatsQuery(ctx, &query)
	if err != nil {
		t.Error(err)
	}
	if len(res.Entries) == 0 {
		t.Errorf("Unexpected empty answer from server")
	}
}

func TestStatsQueryServiceIdent(t *testing.T) {
	query := stats.StatsRequest{}
	query.Discriminator = "node"
	query.StatsCpu = true
	query.StatsMem = true
	query.StatsIo = true
	query.StatsNet = true
	query.Period = "5m"
	query.FilterServiceIdent = "amp"
	res, err := statsClient.StatsQuery(ctx, &query)
	if err != nil {
		t.Error(err)
	}
	if len(res.Entries) == 0 {
		t.Errorf("Unexpected empty answer from server")
	}
	for _, result := range res.Entries {
		if !strings.HasPrefix(result.ServiceName, "amp") {
			t.Errorf("Unexpected service selected: %s\n", result.ServiceName)
		}
	}
}
