package pstats

import "golang.org/x/net/context"

// Stats structure to implement StatsServer interface
type PrometheusStats struct {
}

var sourceMap = make(map[string]*PrometheusSource)

//LoadSources load all prometheus sources
func LoadSources() {
	sourceMap["prometheus"] = newSource("prometheus", "prometheus", "9090")
	sourceMap["node"] = newSource("node", "node_exporter", "9100")
	sourceMap["etcd"] = newSource("etcd", "etcd", "2379")
	sourceMap["haproxy"] = newSource("haproxy", "haproxy_exporter", "9101")
	sourceMap["nats"] = newSource("nats", "nats_exporter", "7777")
	for _, source := range sourceMap {
		source.load()
	}
}

// ReloadSources load all prometheus sources
func (s *PrometheusStats) ReloadSources(ctx context.Context, req *LoadSourcesRequest) (*LoadSourcesReply, error) {
	LoadSources()
	return &LoadSourcesReply{}, nil
}

// PrometheusStats execute any prometheus stats
func (s *PrometheusStats) PrometheusStats(ctx context.Context, req *PrometheusStatsRequest) (*PrometheusStatsReply, error) {
	return nil, nil
}
