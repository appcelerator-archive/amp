package stats

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/appcelerator/amp/pkg/elasticsearch"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
)

// Stats structure to implement StatsServer interface
type Stats struct {
	ES *elasticsearch.Elasticsearch
}

const (
	esIndex = "ampbeat-*"
	esType  = "metrics"
)

// StatsQuery extracts stat information according to StatsRequest
func (s *Stats) StatsQuery(ctx context.Context, req *StatsRequest) (*StatsReply, error) {
	fmt.Printf("request: %+v\n", req)
	if err := s.validatePeriod(req.Period); err != nil {
		return nil, err
	}
	if err := s.validateTimeGroup(req.TimeGroup); err != nil {
		return nil, err
	}
	if err := s.ES.Connect(); err != nil {
		return nil, errors.New("unable to connect to elasticsearch service")
	}
	if req.TimeGroup == "" {
		return s.statsCurrentQuery(ctx, req)
	}
	if req.Group == "" {
		return s.statsHistoricQuery(ctx, req)
	}
	return s.statsGroupedHistoricQuery(ctx, req)
}

// execute a current stats reauest
func (s *Stats) statsCurrentQuery(ctx context.Context, req *StatsRequest) (*StatsReply, error) {
	boolQuery := s.createBoolQuery(req, req.Period)
	agg := s.createTermAggreggation(req)

	result, err := s.ES.GetClient().Search().
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
		return nil, errors.New("Request error 'group' not found")
	}
	ret := &StatsReply{}
	for _, bucket := range ranges.Buckets {
		entry := &MetricsEntry{Group: bucket.Key.(string)}
		if err := s.updateEntry(bucket, req, entry); err != nil {
			return nil, err
		}
		ret.Entries = append(ret.Entries, entry)
	}
	return ret, nil
}

// execute a historic stats request
func (s *Stats) statsHistoricQuery(ctx context.Context, req *StatsRequest) (*StatsReply, error) {
	if req.Period == "" {
		return nil, fmt.Errorf("Historical statistics (using --time-group option) should set --period option explicitelly")
	}
	boolQuery := s.createBoolQuery(req, req.Period)
	agg := s.createHistoAggreggation(req)

	result, err := s.ES.GetClient().Search().
		Index(esIndex).
		Query(boolQuery).
		Aggregation("histo", agg).
		Size(0).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	if result.Hits.TotalHits == 0 {
		return &StatsReply{}, nil
	}
	ranges, ok := result.Aggregations.Terms("histo")
	if !ok {
		return nil, errors.New("Request error 'histo' not found")
	}
	ret := &StatsReply{}
	for _, bucket := range ranges.Buckets {
		entry := &MetricsEntry{Group: (*bucket.KeyAsString)[0:19]}
		if err := s.updateEntry(bucket, req, entry); err != nil {
			return nil, err
		}
		ret.Entries = append(ret.Entries, entry)
	}
	return ret, nil
}

// execute a historic stats request grouped by object name
func (s *Stats) statsGroupedHistoricQuery(ctx context.Context, req *StatsRequest) (*StatsReply, error) {
	if req.Period == "" {
		return nil, fmt.Errorf("Historical statistics (using --time-group option) should set --period option explicitelly")
	}
	boolQuery := s.createBoolQuery(req, req.Period)
	agg := elastic.NewDateHistogramAggregation().Field("@timestamp").Interval(req.TimeGroup)
	aggTerm := s.createTermAggreggation(req)
	agg = agg.SubAggregation("groupByName", aggTerm)

	result, err := s.ES.GetClient().Search().
		Index(esIndex).
		Query(boolQuery).
		Aggregation("histo", agg).
		Size(0).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	if result.Hits.TotalHits == 0 {
		return &StatsReply{}, nil
	}
	ranges, ok := result.Aggregations.Terms("histo")
	if !ok {
		return nil, errors.New("Request error 'histo' not found")
	}
	ret := &StatsReply{}
	for _, bucket1 := range ranges.Buckets {
		ranges2, ok := bucket1.Aggregations.Terms("groupByName")
		if !ok {
			return nil, errors.New("Request error 'groupByName' not found")
		}
		for _, bucket2 := range ranges2.Buckets {
			entry := &MetricsEntry{Group: (*bucket1.KeyAsString)[0:19]}
			entry.Sgroup = bucket2.Key.(string)
			if err := s.updateEntry(bucket2, req, entry); err != nil {
				return nil, err
			}
			ret.Entries = append(ret.Entries, entry)
		}
	}
	return ret, nil
}

// Create the sub-query taking in account all filters and time range
func (s *Stats) createBoolQuery(req *StatsRequest, period string) *elastic.BoolQuery {
	filters := []*elastic.WildcardQuery{}
	if req.FilterContainerId != "" {
		filters = append(filters, elastic.NewWildcardQuery("container_id", getWildcardValue(req.FilterContainerId)))
	}
	if req.FilterContainerName != "" {
		filters = append(filters, elastic.NewWildcardQuery("container_name", getWildcardValue(req.FilterContainerName)))
	}
	if req.FilterContainerShortName != "" {
		filters = append(filters, elastic.NewWildcardQuery("container_short_name", getWildcardValue(req.FilterContainerShortName)))
	}
	if req.FilterContainerState != "" {
		filters = append(filters, elastic.NewWildcardQuery("container_state", getWildcardValue(req.FilterContainerState)))
	}
	if req.FilterServiceName != "" {
		filters = append(filters, elastic.NewWildcardQuery("service_name", getWildcardValue(req.FilterServiceName)))
	}
	if req.FilterServiceId != "" {
		filters = append(filters, elastic.NewWildcardQuery("service_id", getWildcardValue(req.FilterServiceId)))
	}
	if req.FilterStackName != "" {
		filters = append(filters, elastic.NewWildcardQuery("stack_name", getWildcardValue(req.FilterStackName)))
	}
	if req.FilterTaskId != "" {
		filters = append(filters, elastic.NewWildcardQuery("task_id", getWildcardValue(req.FilterTaskId)))
	}
	if req.FilterNodeId != "" {
		filters = append(filters, elastic.NewWildcardQuery("node_id", getWildcardValue(req.FilterNodeId)))
	}
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(elastic.NewRangeQuery("@timestamp").Gte(period), elastic.NewTermsQuery("type", esType))
	for _, query := range filters {
		boolQuery.Must(query)
	}
	return boolQuery
}

func getWildcardValue(val string) string {
	return fmt.Sprintf("%s*", val)
}

// create the aggregation query on the main group (container, service, stacks. ...) and each sub aggregations related to the metrics
func (s *Stats) createTermAggreggation(req *StatsRequest) *elastic.TermsAggregation {
	agg := elastic.NewTermsAggregation().Field(req.Group).Size(100).OrderByTermAsc()
	agg = agg.SubAggregation("countContainer", elastic.NewCardinalityAggregation().Field("container_id"))
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

// create the aggregation query on the main group (container, service, stacks. ...) and each sub aggregations related to the metrics
func (s *Stats) createHistoAggreggation(req *StatsRequest) *elastic.DateHistogramAggregation {
	agg := elastic.NewDateHistogramAggregation().Field("@timestamp").Interval(req.TimeGroup)
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

// extract float value and return error if not exist
func (s *Stats) getInt64CountValue(bucket *elastic.AggregationBucketKeyItem, name string) (float64, error) {
	avg, found := bucket.Cardinality(name)
	if !found {
		return 0, fmt.Errorf("Request error '%s' not found", name)
	}
	value := avg.Value
	if value == nil {
		return 0, nil
	}
	return *value, nil
}

// extract float value and return error if not exist
func (s *Stats) getFloat64AvgValue(bucket *elastic.AggregationBucketKeyItem, name string, nbc float64) (float64, error) {
	avg, found := bucket.Avg(name)
	if !found {
		return 0, fmt.Errorf("Request error '%s' not found", name)
	}
	value := avg.Value
	if value == nil {
		return 0, nil
	}
	return (*value) * nbc, nil
}

// extract int value and return error if not exist
func (s *Stats) getInt64AvgValue(bucket *elastic.AggregationBucketKeyItem, name string, nbc float64) (int64, error) {
	val, err := s.getFloat64AvgValue(bucket, name, nbc)
	if err != nil {
		return 0, err
	}
	return int64(val), nil
}

// update the entry with all the metrics values get by the reauest answer
func (s *Stats) updateEntry(bucket *elastic.AggregationBucketKeyItem, req *StatsRequest, entry *MetricsEntry) error {
	var err error
	var nbc float64 = 1
	if !req.Avg {
		if nbc, err = s.getInt64CountValue(bucket, "countContainer"); err != nil {
			return err
		}
	}
	if req.StatsCpu {
		entry.Cpu = &MetricsCPUEntry{}
		if entry.Cpu.TotalUsage, err = s.getFloat64AvgValue(bucket, "avgCPU", nbc); err != nil {
			return err
		}
		if entry.Cpu.UsageInKernelMode, err = s.getFloat64AvgValue(bucket, "avgCPUKernel", nbc); err != nil {
			return err
		}
		if entry.Cpu.UsageInUserMode, err = s.getFloat64AvgValue(bucket, "avgCPUUser", nbc); err != nil {
			return err
		}
	}
	if req.StatsMem {
		entry.Mem = &MetricsMemEntry{}
		if entry.Mem.Failcnt, err = s.getInt64AvgValue(bucket, "avgMemFailcnt", nbc); err != nil {
			return err
		}
		if entry.Mem.Limit, err = s.getInt64AvgValue(bucket, "avgMemLimit", nbc); err != nil {
			return err
		}
		if entry.Mem.Maxusage, err = s.getInt64AvgValue(bucket, "avgMemMaxUsage", nbc); err != nil {
			return err
		}
		if entry.Mem.Usage, err = s.getInt64AvgValue(bucket, "avgMemUsage", nbc); err != nil {
			return err
		}
		if entry.Mem.UsageP, err = s.getFloat64AvgValue(bucket, "avgMemUsageP", nbc); err != nil {
			return err
		}
	}
	if req.StatsNet {
		entry.Net = &MetricsNetEntry{}
		if entry.Net.TotalBytes, err = s.getInt64AvgValue(bucket, "avgTotalBytes", nbc); err != nil {
			return err
		}
		if entry.Net.RxBytes, err = s.getInt64AvgValue(bucket, "avgRxBytes", nbc); err != nil {
			return err
		}
		if entry.Net.RxDropped, err = s.getInt64AvgValue(bucket, "avgRxDropped", nbc); err != nil {
			return err
		}
		if entry.Net.RxErrors, err = s.getInt64AvgValue(bucket, "avgRxErrors", nbc); err != nil {
			return err
		}
		if entry.Net.RxPackets, err = s.getInt64AvgValue(bucket, "avgRxPackets", nbc); err != nil {
			return err
		}
		if entry.Net.TxBytes, err = s.getInt64AvgValue(bucket, "avgTxBytes", nbc); err != nil {
			return err
		}
		if entry.Net.TxDropped, err = s.getInt64AvgValue(bucket, "avgTxDropped", nbc); err != nil {
			return err
		}
		if entry.Net.TxErrors, err = s.getInt64AvgValue(bucket, "avgTxErrors", nbc); err != nil {
			return err
		}
		if entry.Net.TxPackets, err = s.getInt64AvgValue(bucket, "avgTxPackets", nbc); err != nil {
			return err
		}
	}
	if req.StatsIo {
		entry.Io = &MetricsIOEntry{}
		if entry.Io.Total, err = s.getInt64AvgValue(bucket, "avgIOTotal", nbc); err != nil {
			return err
		}
		if entry.Io.Read, err = s.getInt64AvgValue(bucket, "avgIORead", nbc); err != nil {
			return err
		}
		if entry.Io.Write, err = s.getInt64AvgValue(bucket, "avgIOWrite", nbc); err != nil {
			return err
		}
	}
	return nil
}

func (s *Stats) validatePeriod(rg string) error {
	if rg == "" {
		return nil
	}
	if !strings.HasPrefix(rg, "now-") {
		return fmt.Errorf("period should start y 'now-': %s", rg)
	}
	last := rg[len(rg)-1:]
	if last != "y" && last != "M" && last != "w" && last != "d" && last != "h" && last != "m" && last != "s" {
		return fmt.Errorf("time-group last digit should be in [y,M,w,d,h,m,s]: %s", rg)
	}
	mid := rg[4 : len(rg)-1]
	if _, err := strconv.Atoi(mid); err != nil {
		return fmt.Errorf("period digits between 'now-' and last digit are not numeric: %s", rg)
	}
	return nil
}

func (s *Stats) validateTimeGroup(rg string) error {
	if rg == "" {
		return nil
	}
	last := rg[len(rg)-1:]
	if last != "y" && last != "M" && last != "w" && last != "d" && last != "h" && last != "m" && last != "s" {
		return fmt.Errorf("time-group last digit should be in [y,M,w,d,h,m,s]")
	}
	mid := rg[0 : len(rg)-1]
	num, err := strconv.Atoi(mid)
	if err != nil {
		return fmt.Errorf("the time-group doesn't start by a number")
	}
	if last == "s" && num < 3 {
		return fmt.Errorf("to short time-group, it should be upper than 2s")
	}
	return nil
}
