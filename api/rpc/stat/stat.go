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
	idFieldName := getIDNameFields(req)
	cpuFields := "usage_in_kernelmode, usage_in_usermode, usage_system, usage_total"
	query := buildQueryString(req, cpuFields, idFieldName, "cpu")
	fmt.Println("query: ", query)
	res, err := s.Influx.Query(query)
	if err != nil {
		return nil, err
	}

	cpuReply := CPUReply{}
	if len(res.Results[0].Series) == 0 {
		return nil, errors.New("No result found")
	}
	list := res.Results[0].Series[0].Values
	cpuReply.Entries = make([]*CPUEntry, len(list))
	for i, row := range list {
		entry := CPUEntry{
			Time:           getTimeFieldValue(row[0]),
			Datacenter:     getStringFieldValue(row[1]),
			Host:           getStringFieldValue(row[2]),
			ContainerId:    getStringFieldValue(row[3]),
			ContainerName:  getStringFieldValue(row[4]),
			ContainerImage: getStringFieldValue(row[5]),
			ServiceId:      getStringFieldValue(row[6]),
			ServiceName:    getStringFieldValue(row[7]),
			TaskId:         getStringFieldValue(row[8]),
			TaskName:       getStringFieldValue(row[9]),
			NodeId:         getStringFieldValue(row[10]),
			UsageKernel:    getNumberFieldValue(row[11]),
			UsageUser:      getNumberFieldValue(row[12]),
			UsageSystem:    getNumberFieldValue(row[13]),
			UsageTotal:     getNumberFieldValue(row[14]),
		}
		cpuReply.Entries[i] = &entry
	}
	return &cpuReply, nil
}

func getStringFieldValue(field interface{}) string {
	if field == nil {
		return ""
	}
	return field.(string)
}

func getNumberFieldValue(field interface{}) string {
	if field == nil {
		return "0"
	}
	return field.(json.Number).String()
}

func getTimeFieldValue(field interface{}) int64 {
	if field == nil {
		return 0
	}
	ret, _ := field.(json.Number).Int64()
	return ret
}

//Return specific field name for influx query concidering StatRequest discriminator
func getIDNameFields(req *StatRequest) string {
	var idFieldName = "\"com.docker.swarm.node.id\""
	if req.Discriminator == "container" {
		idFieldName = "cotainer_id"
	} else if req.Discriminator == "service" {
		idFieldName = "\"com.docker.swarm.service.id\""
	} else if req.Discriminator == "task" {
		idFieldName = "\"com.docker.swarm.task.id\""
	} else {
		req.Discriminator = "node"
	}
	return idFieldName
}

//Compute the influx 'sql' query string concidering StatRequest
func buildQueryString(req *StatRequest, fields string, groupby string, metric string) string {
	var where = ""
	//var limit = " LIMIT 1"
	if req.Since != "" {
		where += fmt.Sprintf(" AND time>='%s'", req.Since)
		//limit = ""
	}
	if req.Until != "" {
		where += fmt.Sprintf(" AND time<='%s'", req.Until)
		//limit = ""
	}
	if req.Period != "" {
		where += fmt.Sprintf(" AND time > now() - %s", req.Period)
		//limit = ""
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
	if where != "" {
		where = "WHERE " + where[5:]
	}
	qSelect := `datacenter, host, container_id, container_name, container_image, "com.docker.swarm.service.id", "com.docker.swarm.service.name", "com.docker.swarm.task.id", "com.docker.swarm.task.name", "com.docker.swarm.node.id"`
	return fmt.Sprintf("SELECT %s,%s FROM docker_container_%s %s GROUP BY %s ORDER BY time DESC", qSelect, fields, metric, where, groupby)
}
