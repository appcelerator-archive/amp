package stats

import (
	"fmt"
	"log"

	elastic "gopkg.in/olivere/elastic.v3"

	"github.com/appcelerator/amp/data/elasticsearch"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/pkg/config"
	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// Stats structure to implement StatsServer interface
type Stats struct {
	ElasticsearchURL string
	EsConnected      bool
	Es               *elasticsearch.Elasticsearch
	Store            storage.Interface
	Docker           *dockerClient.Client
	Client           *elastic.Client
}

const (
	esIndex                = "ampbeat-*"
	esType                 = "metrics"
	discriminatorContainer = "container"
	discriminatorService   = "service"
	discriminatorNode      = "node"
	discriminatorTask      = "task"
	metricsCPU             = "cpu"
	metricsMem             = "mem"
	metricsNet             = "net"
	metricsIO              = "io"
)

func (s *Stats) isElasticsearch(ctx context.Context) bool {
	if !s.doesElasticsearchServiceExist(ctx) {
		s.EsConnected = false
		return false
	}
	if !s.EsConnected {
		log.Println("Connecting to elasticsearch at", s.ElasticsearchURL)
		if err := s.Es.Connect(s.ElasticsearchURL, amp.DefaultTimeout); err != nil {
			log.Printf("unable to connect to elasticsearch at %s: %v", s.ElasticsearchURL, err)
			return false
		}
		s.EsConnected = true
		log.Println("Connected to elasticsearch at", s.ElasticsearchURL)
	}
	client := s.Es.GetClient()
	if client.IsRunning() {
		s.Client = client
		return true
	}
	return false
}

func (s *Stats) doesElasticsearchServiceExist(ctx context.Context) bool {
	list, err := s.Docker.ServiceList(ctx, types.ServiceListOptions{
	//Filter: filter,
	})
	if err != nil || len(list) == 0 {
		return false
	}
	for _, serv := range list {
		if serv.Spec.Annotations.Name == "monitoring_elasticsearch" {
			return true
		}
	}
	return false
}

// StatsQuery extracts stat information according to StatsRequest
func (s *Stats) StatsQuery(ctx context.Context, req *StatsRequest) (*StatsReply, error) {
	if !s.isElasticsearch(ctx) {
		return nil, fmt.Errorf("the monitoring_elasticsearch service is not running, please start stack 'monitoring'")
	}
	if req.TimeGroup == "" && req.Period == "" && req.Since == "" && req.Until == "" {
		return s.statsCurrentQuery(ctx, req)
	}
	return s.statsHistoricQuery(ctx, req)
}

func (s *Stats) createBoolQuery(req *StatsRequest, period string) *elastic.BoolQuery {
	filters := []*elastic.TermsQuery{elastic.NewTermsQuery("type", esType)}
	if req.FilterContainerId != "" {
		filters = append(filters, elastic.NewTermsQuery("container_id", req.FilterContainerId))
	}
	if req.FilterContainerName != "" {
		filters = append(filters, elastic.NewTermsQuery("container_name", req.FilterContainerName))
	}
	if req.FilterContainerShortName != "" {
		filters = append(filters, elastic.NewTermsQuery("container_short_name", req.FilterContainerShortName))
	}
	if req.FilterServiceName != "" {
		filters = append(filters, elastic.NewTermsQuery("service_name", req.FilterServiceName))
	}
	if req.FilterServiceId != "" {
		filters = append(filters, elastic.NewTermsQuery("service_id", req.FilterServiceId))
	}
	if req.FilterTaskId != "" {
		filters = append(filters, elastic.NewTermsQuery("task_id", req.FilterTaskId))
	}
	if req.FilterStackName != "" {
		filters = append(filters, elastic.NewTermsQuery("stack_name", req.FilterStackName))
	}
	if req.FilterNodeId != "" {
		filters = append(filters, elastic.NewTermsQuery("node_id", req.FilterNodeId))
	}
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(elastic.NewRangeQuery("@timestamp").Gte(period))
	for _, query := range filters {
		boolQuery.Must(query)
	}
	return boolQuery
}

func (s *Stats) createTermAggreggation(req *StatsRequest) *elastic.TermsAggregation {
	agg := elastic.NewTermsAggregation().Field(req.Group).Size(100).OrderByTermAsc()
	if req.StatsCpu {
		agg = agg.SubAggregation("avgCPU", elastic.NewAvgAggregation().Field("cpu.total_usage"))
		agg = agg.SubAggregation("avgCPUKernel", elastic.NewAvgAggregation().Field("cpu.usage_in_kernel_mode"))
		agg = agg.SubAggregation("avgCPUUser", elastic.NewAvgAggregation().Field("cpu.usage_in_user_mode"))
	}
	if req.StatsMem {
		agg = agg.SubAggregation("avgMemFailcnt", elastic.NewAvgAggregation().Field("mem.fail_count"))
		agg = agg.SubAggregation("avgMemLimit", elastic.NewAvgAggregation().Field("mem.limit"))
		agg = agg.SubAggregation("avgMemMaxUsage", elastic.NewAvgAggregation().Field("mem.max_usage"))
		agg = agg.SubAggregation("avgMemUsage", elastic.NewAvgAggregation().Field("mem.usage"))
		agg = agg.SubAggregation("avgMemUsageP", elastic.NewAvgAggregation().Field("mem.usage_pct"))
	}
	if req.StatsNet {
		agg = agg.SubAggregation("avgTotalBytes", elastic.NewAvgAggregation().Field("net.total_bytes"))
		agg = agg.SubAggregation("avgRxBytes", elastic.NewAvgAggregation().Field("net.rx_bytes"))
		agg = agg.SubAggregation("avgRxDropped", elastic.NewAvgAggregation().Field("net.rx_dropped"))
		agg = agg.SubAggregation("avgRxErrors", elastic.NewAvgAggregation().Field("net.rx_errors"))
		agg = agg.SubAggregation("avgRxPackets", elastic.NewAvgAggregation().Field("net.rx_packets"))
		agg = agg.SubAggregation("avgTxBytes", elastic.NewAvgAggregation().Field("net.tx_bytes"))
		agg = agg.SubAggregation("avgTxDropped", elastic.NewAvgAggregation().Field("net.tx_dropped"))
		agg = agg.SubAggregation("avgTxErrors", elastic.NewAvgAggregation().Field("net.tx_errors"))
		agg = agg.SubAggregation("avgTxPackets", elastic.NewAvgAggregation().Field("net.tx_packets"))
	}
	if req.StatsIo {
		agg = agg.SubAggregation("avgIOTotal", elastic.NewAvgAggregation().Field("io.total"))
		agg = agg.SubAggregation("avgIORead", elastic.NewAvgAggregation().Field("io.read"))
		agg = agg.SubAggregation("avgIOWrite", elastic.NewAvgAggregation().Field("io.write"))
	}
	return agg
}

func (s *Stats) getFloat64AvgValue(bucket *elastic.AggregationBucketKeyItem, name string) (float64, error) {
	avg, found := bucket.Avg(name)
	if !found {
		return 0, fmt.Errorf("Request error '%s' not found", name)
	}
	value := avg.Value
	if value == nil {
		return 0, nil
	}
	return *value, nil
}

func (s *Stats) getInt64AvgValue(bucket *elastic.AggregationBucketKeyItem, name string) (int64, error) {
	val, err := s.getFloat64AvgValue(bucket, name)
	if err != nil {
		return 0, err
	}
	return int64(val), nil
}

func (s *Stats) statsCurrentQuery(ctx context.Context, req *StatsRequest) (*StatsReply, error) {
	boolQuery := s.createBoolQuery(req, "now-10s")
	agg := s.createTermAggreggation(req)

	result, err := s.Client.Search().
		Index(esIndex).
		Query(boolQuery).
		Aggregation("group", agg).
		Size(0).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	if result.Hits.TotalHits == 0 {
		return &StatsReply{}, nil
	}
	ranges, ok := result.Aggregations.Terms("group")
	if !ok {
		return nil, fmt.Errorf("Request error 'group' not found")
	}
	ret := &StatsReply{}
	for _, bucket := range ranges.Buckets {
		entry := &MetricsEntry{Group: bucket.Key.(string)}
		if req.StatsCpu {
			if entry.Cpu.TotalUsage, err = s.getFloat64AvgValue(bucket, "avgCPU"); err != nil {
				return nil, err
			}
			if entry.Cpu.UsageInKernelMode, err = s.getFloat64AvgValue(bucket, "avgCPUKernel"); err != nil {
				return nil, err
			}
			if entry.Cpu.UsageInUserMode, err = s.getFloat64AvgValue(bucket, "avgCPUUser"); err != nil {
				return nil, err
			}
		}
		if req.StatsMem {
			if entry.Mem.Failcnt, err = s.getInt64AvgValue(bucket, "avgMemFailcnt"); err != nil {
				return nil, err
			}
			if entry.Mem.Limit, err = s.getInt64AvgValue(bucket, "avgMemLimit"); err != nil {
				return nil, err
			}
			if entry.Mem.Maxusage, err = s.getInt64AvgValue(bucket, "avgMemMaxUsage"); err != nil {
				return nil, err
			}
			if entry.Mem.Usage, err = s.getInt64AvgValue(bucket, "avgMemUsage"); err != nil {
				return nil, err
			}
			if entry.Mem.UsageP, err = s.getFloat64AvgValue(bucket, "avgMemUsageP"); err != nil {
				return nil, err
			}
		}
		if req.StatsNet {
			if entry.Net.TotalBytes, err = s.getInt64AvgValue(bucket, "avgTotalBytes"); err != nil {
				return nil, err
			}
			if entry.Net.RxBytes, err = s.getInt64AvgValue(bucket, "avgRxBytes"); err != nil {
				return nil, err
			}
			if entry.Net.RxDropped, err = s.getInt64AvgValue(bucket, "avgRxDropped"); err != nil {
				return nil, err
			}
			if entry.Net.RxErrors, err = s.getInt64AvgValue(bucket, "avgRxErrors"); err != nil {
				return nil, err
			}
			if entry.Net.RxPackets, err = s.getInt64AvgValue(bucket, "avgRxPackets"); err != nil {
				return nil, err
			}
			if entry.Net.TxBytes, err = s.getInt64AvgValue(bucket, "avgTxBytes"); err != nil {
				return nil, err
			}
			if entry.Net.TxDropped, err = s.getInt64AvgValue(bucket, "avgTxDropped"); err != nil {
				return nil, err
			}
			if entry.Net.TxErrors, err = s.getInt64AvgValue(bucket, "avgTxErrors"); err != nil {
				return nil, err
			}
			if entry.Net.TxPackets, err = s.getInt64AvgValue(bucket, "avgTxPackets"); err != nil {
				return nil, err
			}
		}
		if req.StatsIo {
			if entry.Io.Total, err = s.getInt64AvgValue(bucket, "avgIOTotal"); err != nil {
				return nil, err
			}
			if entry.Io.Read, err = s.getInt64AvgValue(bucket, "avgIORead"); err != nil {
				return nil, err
			}
			if entry.Io.Write, err = s.getInt64AvgValue(bucket, "avgIOWrite"); err != nil {
				return nil, err
			}
		}
		ret.Entries = append(ret.Entries, entry)
	}
	return ret, nil
}

func (s *Stats) getFloat64HistoAvgValue(bucket *elastic.AggregationBucketHistogramItem, name string) (float64, error) {
	avg, found := bucket.Avg(name)
	if !found {
		return 0, fmt.Errorf("Request error '%s' not found", name)
	}
	value := avg.Value
	if value == nil {
		return 0, nil
	}
	return *value, nil
}

func (s *Stats) getInt64HistoAvgValue(bucket *elastic.AggregationBucketHistogramItem, name string) (int64, error) {
	val, err := s.getFloat64HistoAvgValue(bucket, name)
	if err != nil {
		return 0, err
	}
	return int64(val), nil
}

func (s *Stats) statsHistoricQuery(ctx context.Context, req *StatsRequest) (*StatsReply, error) {
	boolQuery := s.createBoolQuery(req, req.Period)
	agg := s.createTermAggreggation(req)
	histoAgg := elastic.NewDateHistogramAggregation().Field("@timestamp").Interval("30s")
	agg = agg.SubAggregation("histo", histoAgg)

	result, err := s.Client.Search().
		Index(esIndex).
		Query(boolQuery).
		Aggregation("group", agg).
		Size(0).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	if result.Hits.TotalHits == 0 {
		return &StatsReply{}, nil
	}
	ranges, ok := result.Aggregations.Terms("group")
	if !ok {
		return nil, fmt.Errorf("Request error 'group' not found")
	}
	ret := &StatsReply{}
	for _, buck := range ranges.Buckets {
		var found bool
		histo, found := buck.Histogram("histo")
		if !found {
			return nil, fmt.Errorf("Request error 'histo' not found")
		}
		for _, bucket := range histo.Buckets {
			entry := &MetricsEntry{Group: *bucket.KeyAsString}
			if req.StatsCpu {
				if entry.Cpu.TotalUsage, err = s.getFloat64HistoAvgValue(bucket, "avgCPU"); err != nil {
					return nil, err
				}
				if entry.Cpu.UsageInKernelMode, err = s.getFloat64HistoAvgValue(bucket, "avgCPUKernel"); err != nil {
					return nil, err
				}
				if entry.Cpu.UsageInUserMode, err = s.getFloat64HistoAvgValue(bucket, "avgCPUUser"); err != nil {
					return nil, err
				}
			}
			if req.StatsMem {
				if entry.Mem.Failcnt, err = s.getInt64HistoAvgValue(bucket, "avgMemFailcnt"); err != nil {
					return nil, err
				}
				if entry.Mem.Limit, err = s.getInt64HistoAvgValue(bucket, "avgMemLimit"); err != nil {
					return nil, err
				}
				if entry.Mem.Maxusage, err = s.getInt64HistoAvgValue(bucket, "avgMemMaxUsage"); err != nil {
					return nil, err
				}
				if entry.Mem.Usage, err = s.getInt64HistoAvgValue(bucket, "avgMemUsage"); err != nil {
					return nil, err
				}
				if entry.Mem.UsageP, err = s.getFloat64HistoAvgValue(bucket, "avgMemUsageP"); err != nil {
					return nil, err
				}
			}
			if req.StatsNet {
				if entry.Net.TotalBytes, err = s.getInt64HistoAvgValue(bucket, "avgTotalBytes"); err != nil {
					return nil, err
				}
				if entry.Net.RxBytes, err = s.getInt64HistoAvgValue(bucket, "avgRxBytes"); err != nil {
					return nil, err
				}
				if entry.Net.RxDropped, err = s.getInt64HistoAvgValue(bucket, "avgRxDropped"); err != nil {
					return nil, err
				}
				if entry.Net.RxErrors, err = s.getInt64HistoAvgValue(bucket, "avgRxErrors"); err != nil {
					return nil, err
				}
				if entry.Net.RxPackets, err = s.getInt64HistoAvgValue(bucket, "avgRxPackets"); err != nil {
					return nil, err
				}
				if entry.Net.TxBytes, err = s.getInt64HistoAvgValue(bucket, "avgTxBytes"); err != nil {
					return nil, err
				}
				if entry.Net.TxDropped, err = s.getInt64HistoAvgValue(bucket, "avgTxDropped"); err != nil {
					return nil, err
				}
				if entry.Net.TxErrors, err = s.getInt64HistoAvgValue(bucket, "avgTxErrors"); err != nil {
					return nil, err
				}
				if entry.Net.TxPackets, err = s.getInt64HistoAvgValue(bucket, "avgTxPackets"); err != nil {
					return nil, err
				}
			}
			if req.StatsIo {
				if entry.Io.Total, err = s.getInt64HistoAvgValue(bucket, "avgIOTotal"); err != nil {
					return nil, err
				}
				if entry.Io.Read, err = s.getInt64HistoAvgValue(bucket, "avgIORead"); err != nil {
					return nil, err
				}
				if entry.Io.Write, err = s.getInt64HistoAvgValue(bucket, "avgIOWrite"); err != nil {
					return nil, err
				}
			}
			ret.Entries = append(ret.Entries, entry)
		}
	}
	return ret, nil
}
