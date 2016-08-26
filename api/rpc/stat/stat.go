package stat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/appcelerator/amp/data/influx"
	"golang.org/x/net/context"
	//"time"
)

//Stat structure to implement StatServer interface
type Stat struct {
	Influx influx.Influx
}

//CPUQuery Extract CPU information according to StatRequest
func (s *Stat) CPUQuery(ctx context.Context, req *StatRequest) (*CPUReply, error) {
	idFieldName := getIDFieldName(req)
	query := s.buildInfluxQuery(req, "usage_system, usage_total, usage_percent", idFieldName, "cpu")
	fmt.Println("Influx query: "+query)	
	res, err := s.Influx.Query(query)
	if err != nil {
		return nil, err
	}
	if len(res.Results[0].Series) == 0 {
		return nil, errors.New("No result found")
	}
	
	cpuReply := CPUReply{}
	if len(res.Results[0].Series) == 0 {
		return nil, errors.New("No result found")
	}
	list := res.Results[0].Series[0].Values
	cpuReply.Entries = make([]*CPUEntry, len(list))
	for i, row := range list {
		entry := CPUEntry{
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
  			UsageSystem: 		s.getNumberFieldValue(row[11]),
  			UsageTotal: 		s.getNumberFieldValue(row[12]),
  			UsagePercent: 		s.getNumberFieldValue(row[13]),
		}
		cpuReply.Entries[i] = &entry
	}
	return s.computeData(req, &cpuReply)
}

func (s *Stat) computeData(req *StatRequest, data *CPUReply) (*CPUReply, error) {
	// aggreggate rows in map per id concidering req (containner_id | service_id | task_id | nodeId)
	resultMap := make(map[string]*CPUEntry)
	for _, row := range data.Entries {
		key := s.getKey(req, row)
		aggr, ok := resultMap[key]
		if !ok {
			resultMap[key]=row
			if row.UsagePercent != 0 {
				row.Cpu = 1
			}
		} else {
			/*
			if aggr.UsagePercent < row.UsagePercent {
				aggr.UsagePercent = row.UsagePercent
			}
			*/
			aggr.UsagePercent += row.UsagePercent
			if row.UsagePercent != 0 {
				aggr.Cpu++
			}
		}
	}
	// create final result using map
	result := CPUReply{}
	result.Entries = make([]*CPUEntry, len(resultMap))
	var ii int32
    	for key := range resultMap {
        	result.Entries[ii] = resultMap[key]
		ii++
        }
        // copmute cpu usage value for each row
        s.computeCPUUsage(&result)
        return &result, nil
}

func (s *Stat) getKey(req *StatRequest, row *CPUEntry) string {
	if req.Discriminator == "container" {
		return row.ContainerId
	} else if req.Discriminator == "service" {
		return row.ServiceId
	} else if req.Discriminator == "task" {
		return row.TaskId
	}
	return row.NodeId
}


func (s *Stat) computeCPUUsage(cpuReply *CPUReply) {
	for _, row := range cpuReply.Entries {
		if row.Cpu != 0 {
			row.Cpu = row.UsagePercent / row.Cpu
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
func getIDFieldName(req *StatRequest) string {
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
	return idFieldName
}


//Compute the influx 'sql' query string to retriece meta data concidering StatRequest
func (s *Stat) buildInfluxQuery(req *StatRequest, metricFields, idFieldName string, metric string) string {
	if !s.isHistoricQuery(req) {
		req.Period = "1m"
	}
	where := s.buildWhereStatement(req)
	mfields :=`datacenter, host, container_id, container_name, container_image, "com.docker.swarm.service.id", "com.docker.swarm.service.name", "com.docker.swarm.task.id", "com.docker.swarm.task.name", "com.docker.swarm.node.id"`
	return fmt.Sprintf("SELECT %s,%s FROM docker_container_%s %s", mfields, metricFields, metric, where)
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


