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
	idFieldName, nameFieldName := getIDNameFields(req)
	cpuFields := idFieldName + ", " + nameFieldName + ", " + "usage_in_kernelmode, usage_in_usermode, usage_system, usage_total"
	query := buildQueryString(req, cpuFields, idFieldName, "cpu", req.Limit)
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
		ID := row[1].(string)
		var UsageKernel = ""
		if row[3] != nil {
			UsageKernel = row[3].(json.Number).String()
		}
		var UsageUser = ""
		if row[4] != nil {
			UsageUser = row[4].(json.Number).String()
		}
		var UsageSystem = ""
		if row[5] != nil {
			UsageSystem = row[5].(json.Number).String()
		}
		var UsageTotal = ""
		if row[6] != nil {
			UsageTotal = row[6].(json.Number).String()
		}
		entry := CPUEntry{
			Id:          ID,
			Name:        row[2].(string),
			UsageKernel: UsageKernel,
			UsageUser:   UsageUser,
			UsageSystem: UsageSystem,
			UsageTotal:  UsageTotal,
		}
		if err != nil {
			return nil, err
		}
		cpuReply.Entries[i] = &entry
	}
	return &cpuReply, nil
}

//Return specific field name for influx query concidering StatRequest
func getIDNameFields(req *StatRequest) (string, string) {
	var idFieldName = "\"com.docker.swarm.node.id\""
	var nameFieldName = "host"
	if req.Discriminator == "container" {
		idFieldName = "container_id"
		nameFieldName = "container_name"
	} else if req.Discriminator == "service" {
		idFieldName = "\"com.docker.swarm.service.id\""
		nameFieldName = "\"com.docker.swarm.service.name\""
	} else {
		req.Discriminator = "node"
	}
	return idFieldName, nameFieldName
}

//Compute the influx 'sql' query string concidering StatRequest
func buildQueryString(req *StatRequest, fields string, groupby string, metric string, limit string) string {
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
	fmt.Println("where: " + where)
	if limit != "" {
		limit = " LIMIT " + limit
	}
	return fmt.Sprintf("SELECT %s FROM docker_container_%s %s GROUP BY %s ORDER BY time DESC %s", fields, metric, where, groupby, limit)
}
