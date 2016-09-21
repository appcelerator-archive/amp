package stats_test

import (
	"os"
	"testing"

	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/appcelerator/amp/api/server"
	//"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

const (
	queryVal = "measurements"
)

var (
	ctx    context.Context
	client stats.StatsClient
)

func TestMain(m *testing.M) {
	_, conn := server.StartTestServer()
	client = stats.NewStatsClient(conn)
	ctx = context.Background()
	os.Exit(m.Run())
}

func TestStatQueryService(t *testing.T) {
	query := stats.StatsRequest{}
	query.Discriminator = "service"
	query.StatsCpu = true
	query.StatsMem = true
	query.StatsIo = true
	query.StatsNet = true
	query.Period = "5m"
	query.FilterServiceName = "amp"
	res, err := client.StatsQuery(ctx, &query)
	if err != nil {
		t.Error(err)
	}
	if len(res.Entries) == 0 {
		t.Errorf("Unexpected empty answer from server")
	}
}

func TestStatQueryContainer(t *testing.T) {
	query := stats.StatsRequest{}
	query.Discriminator = "container"
	query.StatsCpu = true
	query.StatsMem = true
	query.StatsIo = true
	query.StatsNet = true
	query.Period = "5m"
	query.FilterContainerName = "amp"
	res, err := client.StatsQuery(ctx, &query)
	if err != nil {
		t.Error(err)
	}
	if len(res.Entries) == 0 {
		t.Errorf("Unexpected empty answer from server")
	}
}

func TestStatQueryTask(t *testing.T) {
	query := stats.StatsRequest{}
	query.Discriminator = "task"
	query.StatsCpu = true
	query.StatsMem = true
	query.StatsIo = true
	query.StatsNet = true
	query.Period = "5m"
	query.FilterTaskName = "amp"
	res, err := client.StatsQuery(ctx, &query)
	if err != nil {
		t.Error(err)
	}
	if len(res.Entries) == 0 {
		t.Errorf("Unexpected empty answer from server")
	}
}

func TestStatQueryNode(t *testing.T) {
	query := stats.StatsRequest{}
	query.Discriminator = "node"
	query.StatsCpu = true
	query.StatsMem = true
	query.StatsIo = true
	query.StatsNet = true
	query.Period = "5m"
	res, err := client.StatsQuery(ctx, &query)
	if err != nil {
		t.Error(err)
	}
	if len(res.Entries) == 0 {
		t.Errorf("Unexpected empty answer from server")
	}
}

func TestStatQueryServiceIdent(t *testing.T) {
	query := stats.StatsRequest{}
	query.Discriminator = "node"
	query.StatsCpu = true
	query.StatsMem = true
	query.StatsIo = true
	query.StatsNet = true
	query.Period = "5m"
	query.FilterServiceIdent = "amp"
	res, err := client.StatsQuery(ctx, &query)
	if err != nil {
		t.Error(err)
	}
	if len(res.Entries) == 0 {
		t.Errorf("Unexpected empty answer from server")
	}
}
