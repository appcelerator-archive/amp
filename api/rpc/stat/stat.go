package stat

import (
	"fmt"
	"github.com/appcelerator/amp/data/influx"
	"golang.org/x/net/context"
	"time"
)

//Stat structure to implement StatServer interface
type Stat struct {
	conn influx.Influx
}

// Config is used to provide specific Parameters for Stats Connection
type Config struct {
	Connstr string
	Dbname  string
	U       string
	P       string
	//TODO add pagination?
}

//New return a new implementation of StateServer
func New(cfg Config) (*Stat, error) {

	c := influx.New(cfg.Connstr, cfg.Dbname, cfg.U, cfg.P)
	err := c.Connect(5 * time.Second)
	return &Stat{conn: c}, err
}

//CPUQuery Extract CPU information according to CPURequest
func (s *Stat) CPUQuery(ctx context.Context, req *CPURequest) (*CPUReply, error) {
	if req.RessourceName == "container" {
		return s.containersCPUQuery(ctx, req)
	} else if req.RessourceName == "service" {
		return s.servicesCPUQuery(ctx, req)
	} else {
		return s.nodesCPUQuery(ctx, req)
	}
}

//Extract CPU information for containers according to CPURequest
func (s *Stat) containersCPUQuery(ctx context.Context, req *CPURequest) (*CPUReply, error) {
	query := s.buildQueryString(req, "container_id")
	res, err := s.conn.Query(query)
	if err != nil {
		return nil, err
	}
	cpuReply := CPUReply{}
	list := res.Results[0].Series[0].Values
	cpuReply.Entries = make([]*CPUEntry, len(list))
	for i, row := range list {
		entry, err := s.containerCPUQuery(ctx, req, row[1].(string))
		if err != nil {
			return nil, err
		}
		cpuReply.Entries[i] = entry
	}
	return &cpuReply, nil
}

//Extract CPU information for one container according to CPURequest
func (s *Stat) containerCPUQuery(ctx context.Context, req *CPURequest, containerID string) (*CPUEntry, error) {
	ql := fmt.Sprintf("SELECT container_name, usage_in_kernelmode, usage_in_usermode, usage_system, usage_total FROM docker_container_cpu WHERE container_id='%s' ORDER BY time DESC LIMIT 1", containerID)
	res, err := s.conn.Query(ql)
	if err != nil {
		return nil, err
	}
	data := res.Results[0].Series[0].Values[0]
	entry := CPUEntry{
		ID:          containerID,
		Name:        data[0].(string),
		UsageKernel: data[1].(int32),
		UsageUser:   data[2].(int32),
		UsageSystem: data[3].(int32),
		UsageTotal:  data[4].(int32),
	}
	return &entry, nil
}

//Extract CPU information for services according to CPURequest
func (s *Stat) servicesCPUQuery(ctx context.Context, req *CPURequest) (*CPUReply, error) {
	query := s.buildQueryString(req, "com.docker.swarm.service.id")
	res, err := s.conn.Query(query)
	if err != nil {
		return nil, err
	}
	cpuReply := CPUReply{}
	list := res.Results[0].Series[0].Values
	cpuReply.Entries = make([]*CPUEntry, len(list))
	for i, row := range list {
		entry, err := s.serviceCPUQuery(ctx, req, row[1].(string))
		if err != nil {
			return nil, err
		}
		cpuReply.Entries[i] = entry
	}
	return &cpuReply, nil
}

//Extract CPU information for one service according to CPURequest
func (s *Stat) serviceCPUQuery(ctx context.Context, req *CPURequest, serviceID string) (*CPUEntry, error) {
	ql := fmt.Sprintf("SELECT com.docker.swarm.service.name, usage_in_kernelmode, usage_in_usermode, usage_system, usage_total FROM docker_container_cpu WHERE com.docker.swarm.service.id='%s' ORDER BY time DESC LIMIT 1", serviceID)
	res, err := s.conn.Query(ql)
	if err != nil {
		return nil, err
	}
	data := res.Results[0].Series[0].Values[0]
	entry := CPUEntry{
		ID:          serviceID,
		Name:        data[0].(string),
		UsageKernel: data[1].(int32),
		UsageUser:   data[2].(int32),
		UsageSystem: data[3].(int32),
		UsageTotal:  data[4].(int32),
	}
	return &entry, nil
}

//Extract CPU information for nodes according to CPURequest
func (s *Stat) nodesCPUQuery(ctx context.Context, req *CPURequest) (*CPUReply, error) {
	query := s.buildQueryString(req, "com.docker.swarm.node.id")
	res, err := s.conn.Query(query)
	if err != nil {
		return nil, err
	}
	cpuReply := CPUReply{}
	list := res.Results[0].Series[0].Values
	cpuReply.Entries = make([]*CPUEntry, len(list))
	for i, row := range list {
		entry, err := s.nodeCPUQuery(ctx, req, row[1].(string))
		if err != nil {
			return nil, err
		}
		cpuReply.Entries[i] = entry
	}
	return &cpuReply, nil
}


//Extract CPU information for one node according to CPURequest
func (s *Stat) nodeCPUQuery(ctx context.Context, req *CPURequest, nodeID string) (*CPUEntry, error) {
	ql := fmt.Sprintf("SELECT host, usage_in_kernelmode, usage_in_usermode, usage_system, usage_total FROM docker_container_cpu WHERE com.docker.swarm.node.id='%s' ORDER BY time DESC LIMIT 1", nodeID)
	res, err := s.conn.Query(ql)
	if err != nil {
		return nil, err
	}
	data := res.Results[0].Series[0].Values[0]
	entry := CPUEntry{
		ID:          nodeID,
		Name:        data[0].(string),
		UsageKernel: data[1].(int32),
		UsageUser:   data[2].(int32),
		UsageSystem: data[3].(int32),
		UsageTotal:  data[4].(int32),
	}
	return &entry, nil
}

func (s *Stat) buildQueryString(req *CPURequest, ressourceFieldName string) string {
	query := fmt.Sprintf("SELECT distinct %s FROM docker_container_cpu WHERE usage_total>0", ressourceFieldName)
	if (req.FilterDatacenter != "") {
		query+= fmt.Sprintf(" AND datacenter='%s'", req.FilterDatacenter)
	}
	if (req.FilterServiceName != "") {
		query+= fmt.Sprintf(" AND com.docker.swarm.service.name='%s'", req.FilterServiceName)
	}
	if (req.FilterServiceID != "") {
		query+= fmt.Sprintf(" AND com.docker.swarm.service.id='%s'", req.FilterServiceName)
	}
	if (req.FilterNodeID != "") {
		query+= fmt.Sprintf(" AND com.docker.swarm.node.id='%s'", req.FilterDatacenter)
	}
	if (req.FilterContainerName != "") {
		query+= fmt.Sprintf(" AND container_name='%s'", req.FilterDatacenter)
	}
	if (req.FilterContainerImage != "") {
		query+= fmt.Sprintf(" AND container_image='%s'", req.FilterDatacenter)
	}
	return query
}
