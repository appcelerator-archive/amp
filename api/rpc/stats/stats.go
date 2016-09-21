package stats

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/appcelerator/amp/data/influx"
	"golang.org/x/net/context"
)

// Stats structure to implement StatsServer interface
type Stats struct {
	Influx influx.Influx
}

const (
	discriminatorContainer = "container"
	discriminatorService   = "service"
	discriminatorNode      = "node"
	discriminatorTask      = "task"
	metricsCPU             = "cpu"
	metricsMem             = "mem"
	metricsNet             = "net"
	metricsIO              = "io"
)

// StatsQuery extracts stat information according to StatsRequest
func (s *Stats) StatsQuery(ctx context.Context, req *StatsRequest) (*StatsReply, error) {
	var metricList [4]*StatsReply
	if req.StatsCpu {
		ret, err := s.statQueryMetric(req, metricsCPU)
		if err != nil {
			return nil, err
		}
		s.addStatsResult(&metricList, ret)
	}
	if req.StatsMem {
		ret, err := s.statQueryMetric(req, metricsMem)
		if err != nil {
			return nil, err
		}
		s.addStatsResult(&metricList, ret)
	}
	if req.StatsIo {
		ret, err := s.statQueryMetric(req, metricsIO)
		if err != nil {
			return nil, err
		}
		s.addStatsResult(&metricList, ret)
	}
	if req.StatsNet {
		ret, err := s.statQueryMetric(req, metricsNet)
		if err != nil {
			return nil, err
		}
		s.addStatsResult(&metricList, ret)
	}
	//fmt.Println(metricList)
	result := s.combineStats(req, &metricList)
	if result == nil || result.Entries == nil {
		ret := &StatsReply{
			Entries: make([]*StatsEntry, 0),
		}
		return ret, nil
	}
	sort.Sort(result)
	return result, nil
}

func (s *Stats) addStatsResult(list *[4]*StatsReply, ret *StatsReply) {
	//debugList("--", list)
	for i := 0; i < 4; i++ {
		if list[i] == nil {
			list[i] = ret
			break
		} else if len(list[i].Entries) < len(ret.Entries) {
			for j := 2; j >= i; j-- {
				list[j+1] = list[j]
			}
			list[i] = ret
			break
		}
	}
}

func (s *Stats) combineStats(req *StatsRequest, list *[4]*StatsReply) *StatsReply {
	finalRet := list[0]
	for i := 1; i < 4; i++ {
		if list[i] != nil {
			ret := list[i]
			for _, frow := range finalRet.Entries {
				for _, row := range ret.Entries {
					if s.isRowsMatch(req, frow, row) {
						s.updateRow(frow, row)
					}
				}
			}
		}
	}
	return finalRet

}

/*
func debugList(mes string, list *[4]*StatsReply) {
  fmt.Println(mes)
  for i := 0 ; i < 4 ; i++ {
    if list[i]==nil {
      fmt.Printf("%d nil\n", i)
    } else {
      fmt.Printf("%d %d\n", i, len(list[i].Entries))
    }
  }
}
*/

func (s *Stats) isRowsMatch(req *StatsRequest, r1 *StatsEntry, r2 *StatsEntry) bool {
	if req.Discriminator == discriminatorContainer {
		if r1.ContainerId == r2.ContainerId {
			return true
		}
		return false
	} else if req.Discriminator == discriminatorService {
		if r1.ServiceId == r2.ServiceId {
			return true
		}
		return false
	} else if req.Discriminator == discriminatorTask {
		if r1.TaskId == r2.TaskId {
			return true
		}
		return false
	}
	if r1.NodeId == r2.NodeId {
		return true
	}
	return false
}

func (s *Stats) updateRow(ref *StatsEntry, row *StatsEntry) {
	if row.Type == metricsCPU {
		ref.Cpu = row.Cpu
	} else if row.Type == metricsMem {
		ref.Mem = row.Mem
		ref.MemUsage = row.MemUsage
		ref.MemLimit = row.MemLimit
	} else if row.Type == metricsIO {
		ref.IoRead = row.IoRead
		ref.IoWrite = row.IoWrite
	} else if row.Type == metricsNet {
		ref.NetTxBytes = row.NetTxBytes
		ref.NetRxBytes = row.NetRxBytes
	}
}

// statsQueryMetric extracts stat information according to StatsRequest for one  metric (cpu | mem | io | net)
func (s *Stats) statQueryMetric(req *StatsRequest, metric string) (*StatsReply, error) {
	idFieldName, metricFields := getMetricFieldsName(req, metric)
	query := s.buildInfluxQuery(req, metricFields, idFieldName, metric)
	fmt.Println("Influx query: " + query)
	res, err := s.Influx.Query(query)
	if err != nil {
		return nil, err
	}
	if len(res.Results[0].Series) == 0 {
		ret := &StatsReply{
			Entries: make([]*StatsEntry, 0),
		}
		return ret, nil
	}
	list := res.Results[0].Series[0].Values
	containerMap := make(map[string]*StatsEntry)
	for _, row := range list {
		entry := StatsEntry{
			Time:           s.getTimeFieldValue(row[0]),
			Datacenter:     s.getStringFieldValue(row[1]),
			Host:           s.getStringFieldValue(row[2]),
			ContainerId:    s.getStringFieldValue(row[3]),
			ContainerName:  s.getStringFieldValue(row[4]),
			ContainerImage: s.getStringFieldValue(row[5]),
			ServiceId:      s.getStringFieldValue(row[6]),
			ServiceName:    s.getStringFieldValue(row[7]),
			TaskId:         s.getStringFieldValue(row[8]),
			TaskName:       s.getStringFieldValue(row[9]),
			NodeId:         s.getStringFieldValue(row[10]),
			SortType:       req.Discriminator,
		}
		entry.Type = metric
		if metric == metricsCPU {
			entry.Cpu = s.getNumberFieldValue(row[11])
		} else if metric == metricsMem {
			entry.Mem = s.getNumberFieldValue(row[11])
			entry.MemUsage = s.getNumberFieldValue(row[12])
			entry.MemLimit = s.getNumberFieldValue(row[13])
		} else if metric == metricsIO {
			entry.IoRead = s.getNumberFieldValue(row[11])
			entry.IoWrite = s.getNumberFieldValue(row[12])
		} else if metric == metricsNet {
			entry.NetTxBytes = s.getNumberFieldValue(row[11])
			entry.NetRxBytes = s.getNumberFieldValue(row[12])
		}
		s.avgInContainerMap(containerMap, &entry)
	}
	return s.addByKeyUsingContainerData(req, containerMap)
}

func (s *Stats) avgInContainerMap(containerMap map[string]*StatsEntry, row *StatsEntry) {
	key := row.ContainerId
	aggr, ok := containerMap[key]
	if !ok {
		containerMap[key] = row
		if row.Cpu != 0 || row.Mem != 0 || row.IoRead != 0 || row.IoWrite != 0 || row.NetTxBytes != 0 || row.NetRxBytes != 0 {
			row.Number = 1
		}
	} else {
		aggr.Cpu += row.Cpu
		aggr.Mem += row.Mem
		aggr.MemUsage += row.MemUsage
		aggr.MemLimit += row.MemLimit
		aggr.IoRead += row.IoRead
		aggr.IoWrite += row.IoWrite
		aggr.NetTxBytes += row.NetTxBytes
		aggr.NetRxBytes += row.NetRxBytes
		if row.Cpu != 0 || row.Mem != 0 {
			aggr.Number++
		}
	}
}

func (s *Stats) addByKeyUsingContainerData(req *StatsRequest, containerMap map[string]*StatsEntry) (*StatsReply, error) {
	// aggreggate rows in map per id concidering req (containner_id | service_id | task_id | nodeId)
	resultMap := make(map[string]*StatsEntry)
	for _, row := range containerMap {
		if row.Number > 0 {
			key := s.getKey(req, row)
			aggr, ok := resultMap[key]
			if !ok {
				aggr = row
				aggr.Cpu = (row.Cpu / row.Number)
				aggr.Mem = (row.Mem / row.Number)
				aggr.MemUsage = (row.MemUsage / row.Number)
				aggr.MemLimit = (row.MemLimit / row.Number)
				aggr.IoRead = (row.IoRead / row.Number)
				aggr.IoWrite = (row.IoWrite / row.Number)
				aggr.NetTxBytes = (row.NetTxBytes / row.Number)
				aggr.NetRxBytes = (row.NetRxBytes / row.Number)
				resultMap[key] = aggr
			} else {
				aggr.Cpu += (row.Cpu / row.Number)
				aggr.Mem += (row.Mem / row.Number)
				aggr.MemUsage += (row.MemUsage / row.Number)
				aggr.MemLimit += (row.MemLimit / row.Number)
				aggr.IoRead += (row.IoRead / row.Number)
				aggr.IoWrite += (row.IoWrite / row.Number)
				aggr.NetTxBytes += (row.NetTxBytes / row.Number)
				aggr.NetRxBytes += (row.NetRxBytes / row.Number)
			}
		}
	}
	// create final result using map
	result := StatsReply{}
	result.Entries = make([]*StatsEntry, len(resultMap))
	var ii int32
	for key, row := range resultMap {
		if key != "" {
			result.Entries[ii] = row
			ii++
		}
	}
	return &result, nil
}

func (s *Stats) getKey(req *StatsRequest, row *StatsEntry) string {
	if !s.isHistoricQuery(req) {
		if req.Discriminator == discriminatorContainer {
			return row.ContainerId
		} else if req.Discriminator == discriminatorService {
			return row.ServiceId
		} else if req.Discriminator == discriminatorTask {
			return row.TaskId
		}
		return row.NodeId
	}
	var period = "m"
	if req.Period != "" {
		period = req.Period[len(req.Period)-1 : len(req.Period)]
	}
	if period == "m" {
		return fmt.Sprintf("%d", row.Time/60)
	} else if period == "h" {
		return fmt.Sprintf("%d", row.Time/3600)
	} else if period == "d" {
		return fmt.Sprintf("%d", row.Time/(3600*24))
	} else if period == "w" {
		return fmt.Sprintf("%d", row.Time/(3600*24*7))
	}
	return fmt.Sprintf("%d", row.Time)
}

func (s *Stats) getStringFieldValue(field interface{}) string {
	if field == nil {
		return ""
	}
	return field.(string)
}

func (s *Stats) getNumberFieldValue(field interface{}) float64 {
	if field == nil {
		return 0
	}
	val, err := field.(json.Number).Float64()
	if err != nil {
		return 0
	}
	return val
}

func (s *Stats) getTimeFieldValue(field interface{}) int64 {
	if field == nil {
		return 0
	}
	val, err := field.(json.Number).Int64()
	if err != nil {
		return 0
	}
	return val
}

func (s *Stats) isHistoricQuery(req *StatsRequest) bool {
	if req.Period != "" || req.Since != "" || req.Until != "" {
		return true
	}
	return false
}

// Return specific field name for influx query concidering StatsRequest discriminator
func getMetricFieldsName(req *StatsRequest, metric string) (string, string) {
	var idFieldName = "\"com.docker.swarm.node.id\""
	if req.Discriminator == discriminatorContainer {
		idFieldName = "container_id"
	} else if req.Discriminator == discriminatorService {
		idFieldName = "\"com.docker.swarm.service.id\""
	} else if req.Discriminator == discriminatorTask {
		idFieldName = "\"com.docker.swarm.task.id\""
	} else {
		req.Discriminator = discriminatorNode
	}
	var fields string
	if metric == metricsCPU {
		fields = "usage_percent"
	} else if metric == metricsMem {
		fields = "usage_percent, usage, max_usage"
	} else if metric == metricsIO {
		fields = "io_serviced_recursive_read, io_serviced_recursive_write"
	} else if metric == metricsNet {
		fields = "rx_bytes, tx_bytes"
	}
	return idFieldName, fields
}

// Compute the influx 'sql' query string to retriece meta data concidering StatsRequest
func (s *Stats) buildInfluxQuery(req *StatsRequest, metricFields, idFieldName string, metric string) string {
	where := s.buildWhereStatsement(req)
	mfields := `datacenter, host, container_id, container_name, container_image, "com.docker.swarm.service.id", "com.docker.swarm.service.name", "com.docker.swarm.task.id", "com.docker.swarm.task.name", "com.docker.swarm.node.id"`
	if metric == "io" {
		metric = "blkio"
	}
	return fmt.Sprintf("SELECT %s, %s FROM docker_container_%s %s", mfields, metricFields, metric, where)
}

// Compute the influx 'sql' WHERE statement concidering the StatsRequest, manage all the filters
func (s *Stats) buildWhereStatsement(req *StatsRequest) string {
	var where = ""
	if req.Since != "" {
		where += fmt.Sprintf(" AND time>='%s'", req.Since)
	}
	if req.Until != "" {
		where += fmt.Sprintf(" AND time<='%s'", req.Until)
	}
	if req.Period != "" {
		where += fmt.Sprintf(" AND time > now() - %s", req.Period)
	} else if !s.isHistoricQuery(req) {
		where += fmt.Sprintf(" AND time > now() - %s", "1m")
	}
	if req.FilterDatacenter != "" {
		where += fmt.Sprintf(" AND datacenter =~ /%s.*/", req.FilterDatacenter)
	}
	if req.FilterHost != "" {
		where += fmt.Sprintf(" AND host =~ /%s.*/", req.FilterHost)
	}
	if req.FilterContainerId != "" {
		where += fmt.Sprintf(" AND container_id =~ /%s.*/", req.FilterContainerId)
	}
	if req.FilterContainerName != "" {
		where += fmt.Sprintf(" AND container_name =~ /%s.*/", req.FilterContainerName)
	}
	if req.FilterContainerImage != "" {
		where += fmt.Sprintf(" AND container_image =~ /%s.*/", req.FilterContainerImage)
	}
	if req.FilterServiceId != "" {
		where += fmt.Sprintf(" AND \"com.docker.swarm.service.id\" =~ /%s.*/", req.FilterServiceId)
	}
	if req.FilterServiceName != "" {
		where += fmt.Sprintf(" AND \"com.docker.swarm.service.name\" =~ /%s.*/", req.FilterServiceName)
	}
	if req.FilterServiceIdent != "" {
		where += fmt.Sprintf(" AND (\"com.docker.swarm.service.name\" =~ /%s.*/ OR \"com.docker.swarm.service.id\" =~ /%s.*/ )", req.FilterServiceIdent, req.FilterServiceIdent)
	}
	if req.FilterTaskId != "" {
		where += fmt.Sprintf(" AND \"com.docker.swarm.task.id\" =~ /%s.*/", req.FilterTaskId)
	}
	if req.FilterTaskName != "" {
		where += fmt.Sprintf(" AND \"com.docker.swarm.task.name\" =~ /%s.*/", req.FilterTaskName)
	}
	if req.FilterNodeId != "" {
		where += fmt.Sprintf(" AND \"com.docker.swarm.node.id\" =~ /%s.*/", req.FilterNodeId)
	}
	if where == "" {
		return ""
	}
	return "WHERE " + where[5:]
}

func getSortKeyValue(row *StatsEntry) string {
	if row.SortType == discriminatorContainer {
		return row.ContainerId
	} else if row.SortType == discriminatorService {
		return row.ServiceId
	} else if row.SortType == discriminatorTask {
		return row.TaskId
	}
	return row.NodeId
}

func (a StatsReply) Len() int {
	return len(a.Entries)
}

func (a StatsReply) Swap(i, j int) {
	a.Entries[i], a.Entries[j] = a.Entries[j], a.Entries[i]
}

func (a StatsReply) Less(i, j int) bool {
	ret := strings.Compare(getSortKeyValue(a.Entries[i]), getSortKeyValue(a.Entries[j]))
	if ret == 0 {
		if a.Entries[i].Time < a.Entries[j].Time {
			return true
		}
		return false
	} else if ret == -1 {
		return true
	}
	return false
}
