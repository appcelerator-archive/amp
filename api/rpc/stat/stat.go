package stat

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sort"
	"github.com/appcelerator/amp/data/influx"
	"golang.org/x/net/context"
	//"time"
)

//Stat structure to implement StatServer interface
type Stat struct {
	Influx influx.Influx
}

//StatQuery Extract stat information according to StatRequest
func (s *Stat) StatQuery(ctx context.Context, req *StatRequest) (*StatReply, error) {
	var metricList [4]*StatReply
	if (req.StatCpu) {
		ret, err := s.statQueryMetric(req, "cpu")
		if err != nil {
			return nil, err
		}
		s.addStatResult(&metricList, ret)
	}
	if (req.StatMem) {
		ret, err := s.statQueryMetric(req, "mem")
		if err != nil {
			return nil, err
		}
		s.addStatResult(&metricList, ret)
	}
	if (req.StatIo) {
		ret, err := s.statQueryMetric(req, "io")
		if err != nil {
			return nil, err
		}
		s.addStatResult(&metricList, ret)
	}
	if (req.StatNet) {
		ret, err := s.statQueryMetric(req, "net")
		if err != nil {
			return nil, err
		}
		s.addStatResult(&metricList, ret)
	}
	//fmt.Println(metricList)
	result := s.combineStat(req, &metricList)
	sort.Sort(result)
	return result, nil

}

func (s *Stat) addStatResult(list *[4]*StatReply, ret *StatReply) {
	//debugList("--", list)
	for i := 0; i < 4; i++ {
		if list[i] == nil {
			list[i] = ret
			break
		} else if len(list[i].Entries) > len(ret.Entries) {
			for j :=2; j>=i ; j-- {
				list[j+1] = list[j]
			}
			list[i] = ret
			break
		}
	}
}

func (s *Stat) combineStat(req *StatRequest, list *[4]*StatReply) *StatReply {
	finalRet := list[0]
	for i := 1 ; i < 4 ; i++ {
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

func debugList(mes string, list *[4]*StatReply) {
	fmt.Println(mes)
	for i := 0 ; i < 4 ; i++ {
		if list[i]==nil {
			fmt.Printf("%d nil\n", i)
		} else {
			fmt.Printf("%d %d\n", i, len(list[i].Entries))
		}
	}
}

func (s *Stat) isRowsMatch(req *StatRequest, r1 *StatEntry, r2 *StatEntry) bool {
	if req.Discriminator == "container" {
		if r1.ContainerId == r2.ContainerId {
			return true
		}
		return false
	} else if req.Discriminator == "service" {
		if r1.ServiceId == r2.ServiceId {
			return true
		}
		return false
	} else if req.Discriminator == "task" {
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

func (s *Stat) updateRow(ref *StatEntry, row *StatEntry) {
	if row.Type == "cpu" {
		ref.Cpu = row.Cpu
	} else if row.Type == "mem" {
		ref.Mem = row.Mem
		ref.MemUsage = row.MemUsage
		ref.MemLimit = row.MemLimit
	} else if row.Type == "io" {
		ref.IoRead = row.IoRead
		ref.IoWrite = row.IoWrite
	} else if row.Type == "net" {
		ref.NetTxBytes = row.NetTxBytes
		ref.NetRxBytes = row.NetRxBytes
	}
}

//StatQuery Extract stat information according to StatRequest for one  metric (cpu | mem | io | net)
func (s *Stat) statQueryMetric(req *StatRequest, metric string) (*StatReply, error) {
	idFieldName, metricFields := getMetricFieldsName(req, metric)
	query := s.buildInfluxQuery(req, metricFields, idFieldName, metric)
	fmt.Println("Influx query: "+query)	
	res, err := s.Influx.Query(query)
	if err != nil {
		return nil, err
	}
	if len(res.Results[0].Series) == 0 {
		return nil, errors.New("No result found")
	}
	
	cpuReply := StatReply{}
	if len(res.Results[0].Series) == 0 {
		return nil, errors.New("No result found")
	}
	list := res.Results[0].Series[0].Values
	cpuReply.Entries = make([]*StatEntry, len(list))
	for i, row := range list {
		entry := StatEntry{
			Time:	 		s.getTimeFieldValue(row[0]),
			Datacenter: 		s.getStringFieldValue(row[1]),
  			Host:			s.getStringFieldValue(row[2]),
  			ContainerId:		s.getStringFieldValue(row[3]),
  			ContainerName:		s.getStringFieldValue(row[4]),
  			ContainerImage:		s.getStringFieldValue(row[5]),
  			ServiceId:		s.getStringFieldValue(row[6]),
  			ServiceName:		s.getStringFieldValue(row[7]),
  			TaskId:			s.getStringFieldValue(row[8]),
  			TaskName:		s.getStringFieldValue(row[9]),
  			NodeId:			s.getStringFieldValue(row[10]),
		}
		entry.Type = metric
		if metric == "cpu" {
			entry.Cpu = s.getNumberFieldValue(row[11])
		} else if metric == "mem" {
			entry.Mem = s.getNumberFieldValue(row[11])
			entry.MemUsage = s.getNumberFieldValue(row[12])
			entry.MemLimit = s.getNumberFieldValue(row[13])
		} else if metric == "io" {
			entry.IoRead = s.getNumberFieldValue(row[11])
			entry.IoWrite = s.getNumberFieldValue(row[12])
		} else if metric == "net" {
			entry.NetTxBytes = s.getNumberFieldValue(row[11])
			entry.NetRxBytes = s.getNumberFieldValue(row[12])
		}

		cpuReply.Entries[i] = &entry
	}
	return s.computeData(req, &cpuReply)
}


func (s *Stat) computeData(req *StatRequest, data *StatReply) (*StatReply, error) {
	// aggreggate rows in map per id concidering req (containner_id | service_id | task_id | nodeId)
	resultMap := make(map[string]*StatEntry)
	for _, row := range data.Entries {
		key := s.getKey(req, row)
		aggr, ok := resultMap[key]
		if !ok {
			resultMap[key]=row
			if row.Cpu != 0 || row.Mem != 0 || row.IoRead !=0 || row.IoWrite != 0 || row.NetTxBytes != 0 || row.NetRxBytes != 0 {
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
	// create final result using map
	result := StatReply{}
	result.Entries = make([]*StatEntry, len(resultMap))
	var ii int32
    	for key := range resultMap {
        	result.Entries[ii] = resultMap[key]
		ii++
        }
        // copmute cpu usage value for each row
        s.computeMetric(&result)
        return &result, nil
}

func (s *Stat) getKey(req *StatRequest, row *StatEntry) string {
	if !s.isHistoricQuery(req) {
		if req.Discriminator == "container" {
			return row.ContainerId
		} else if req.Discriminator == "service" {
			return row.ServiceId
		} else if req.Discriminator == "task" {
			return row.TaskId
		}
		return row.NodeId
	}
	var period = "m"
	if req.Period != "" {
		period = req.Period[len(req.Period) - 1:len(req.Period)]
	}
	if period == "m" {
		return fmt.Sprintf("%d", row.Time / 60)
	} else if period == "h" {
		return fmt.Sprintf("%d", row.Time / 3600)
	} else if period == "d" {
		return fmt.Sprintf("%d", row.Time / (3600 * 24))
	} else if period == "w" {
		return fmt.Sprintf("%d", row.Time / (3600* 24 * 7))
	}
	return fmt.Sprintf("%d", row.Time / 60)
 }


func (s *Stat) computeMetric(cpuReply *StatReply) {
	for _, row := range cpuReply.Entries {
		if row.Cpu != 0 {
			row.Cpu = row.Cpu / row.Number
		}
		if row.Mem != 0 {
			row.Mem = row.Mem / row.Number
			row.MemUsage = row.MemUsage / row.Number
			row.MemLimit = row.MemLimit / row.Number
		}
		if row.IoRead != 0 {
			row.IoRead = row.IoRead / row.Number
		}
		if row.IoWrite != 0 {
			row.IoWrite = row.IoWrite / row.Number
		}
		if row.NetTxBytes != 0 {
			row.NetTxBytes = row.NetTxBytes / row.Number
		}
		if row.NetRxBytes != 0 {
			row.NetRxBytes = row.NetRxBytes / row.Number
		}
	}
}

func (s *Stat) getStringFieldValue(field interface {}) string {
	if field == nil {
		return ""
	}
	return field.(string)
}

func (s *Stat) getNumberFieldValue(field interface {}) float64 {
	if field == nil {
		return 0
	}
	val, err := field.(json.Number).Float64()
	if err != nil {
		return 0
	}
	return val
}

func (s *Stat) getTimeFieldValue(field interface {}) int64 {
	if field == nil {
		return 0
	}
	val, err := field.(json.Number).Int64()
	if err != nil {
		return 0
	}
	return val
}

func (s *Stat) isHistoricQuery(req *StatRequest) bool {
	if req.Period != "" || req.Since != "" || req.Until != "" {
		return true
	}
	return false
}

//Return specific field name for influx query concidering StatRequest discriminator
func getMetricFieldsName(req *StatRequest, metric string) (string,  string) {
	var idFieldName = "\"com.docker.swarm.node.id\""
	if req.Discriminator == "container" {
		idFieldName = "container_id"
	} else if req.Discriminator == "service" {
		idFieldName = "\"com.docker.swarm.service.id\""
	} else if req.Discriminator == "task" {
		idFieldName = "\"com.docker.swarm.task.id\""
	} else {
		req.Discriminator = "node"
	}
	var fields string
	if metric == "cpu" {
		fields = "usage_percent"
	} else if metric == "mem" {
		fields = "usage_percent, usage, max_usage"
	} else if metric == "io" {
		fields = "io_serviced_recursive_read, io_serviced_recursive_write"
	} else if metric == "net" {
		fields = "rx_bytes, tx_bytes"
	}
	return idFieldName, fields
}


//Compute the influx 'sql' query string to retriece meta data concidering StatRequest
func (s *Stat) buildInfluxQuery(req *StatRequest, metricFields, idFieldName string, metric string) string {
	where := s.buildWhereStatement(req)
	mfields :=`datacenter, host, container_id, container_name, container_image, "com.docker.swarm.service.id", "com.docker.swarm.service.name", "com.docker.swarm.task.id", "com.docker.swarm.task.name", "com.docker.swarm.node.id"`
	if metric == "io" {
		metric = "blkio"
	}
	return fmt.Sprintf("SELECT %s, %s FROM docker_container_%s %s", mfields, metricFields, metric, where)
}


//Compute the influx 'sql' WHERE statement concidering the StatRequest, manage all the filters
func (s *Stat) buildWhereStatement(req *StatRequest) string {
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
		where += fmt.Sprintf(" AND datacenter='%s'", req.FilterDatacenter)
	}
	if req.FilterHost != "" {
		where += fmt.Sprintf(" AND host='%s'", req.FilterHost)
	}
	if req.FilterContainerId != "" {
		where += fmt.Sprintf(" AND container_id='%s'", req.FilterContainerId)
	}
	if req.FilterContainerName != "" {
		where += fmt.Sprintf(" AND container_name='%s'", req.FilterContainerName)
	}
	if req.FilterContainerImage != "" {
		where += fmt.Sprintf(" AND container_image='%s'", req.FilterContainerImage)
	}
	if req.FilterServiceId != "" {
		where += fmt.Sprintf(" AND \"com.docker.swarm.service.id\"='%s'", req.FilterServiceId)
	}
	if req.FilterServiceName != "" {
		where += fmt.Sprintf(" AND \"com.docker.swarm.service.name\"='%s'", req.FilterServiceName)
	}
	if req.FilterTaskId != "" {
		where += fmt.Sprintf(" AND \"com.docker.swarm.task.id\"='%s'", req.FilterTaskId)
	}
	if req.FilterTaskName != "" {
		where += fmt.Sprintf(" AND \"com.docker.swarm.task.name\"='%s'", req.FilterTaskName)
	}
	if req.FilterNodeId != "" {
		where += fmt.Sprintf(" AND \"com.docker.swarm.node.id\"='%s'", req.FilterNodeId)
	}
	if where == "" {
		return ""
	}
	return "WHERE " + where[5:]
}


func (a StatReply) Len() int { 
	return len(a.Entries) 
}

func (a StatReply) Swap(i, j int) { 
	a.Entries[i], a.Entries[j] = a.Entries[j], a.Entries[i] 
}

func (a StatReply) Less(i, j int) bool { 
	if ret := strings.Compare(a.Entries[i].ContainerId, a.Entries[j].ContainerId); ret ==-1 {
		return true
	} 
	return false
}