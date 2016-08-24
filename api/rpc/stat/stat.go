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


//CPUQuery Extract CPU information according to StatRequest
func (s *Stat) CPUQuery(ctx context.Context, req *StatRequest) (*CPUReply, error) {

	idFieldName, nameFieldName := getIDNameFields(req)
	query := buildQueryString(req, idFieldName, "cpu")
	res, err := s.conn.Query(query)
	if err != nil {
		return nil, err
	}

	cpuReply := CPUReply{}
	list := res.Results[0].Series[0].Values
	cpuReply.Entries = make([]*CPUEntry, len(list))
	for i, row := range list {
		entry, err := s.ressourceCPUQuery(ctx, req, idFieldName, nameFieldName, row[1].(string))
		if err != nil {
			return nil, err
		}
		cpuReply.Entries[i] = entry
	}
	return &cpuReply, nil
}

//Extract CPU information for one ressource according to StatRequest
func (s *Stat) ressourceCPUQuery(ctx context.Context, req *StatRequest, idFieldName string, nameFieldName, ID string) (*CPUEntry, error) {
	var cpuFields = "usage_in_kernelmode, usage_in_usermode, usage_system, usage_total"
	ql := fmt.Sprintf("SELECT %s, %s FROM docker_container_cpu WHERE %s='%s' ORDER BY time DESC LIMIT 1", nameFieldName, cpuFields, idFieldName, ID)
	res, err := s.conn.Query(ql)
	if err != nil {
		return nil, err
	}
	data := res.Results[0].Series[0].Values[0]
	entry := CPUEntry{
		Id:          ID,
		Name:        data[0].(string),
		UsageKernel: data[1].(int64),
		UsageUser:   data[2].(int64),
		UsageSystem: data[3].(int64),
		UsageTotal:  data[4].(int64),
	}
	return &entry, nil
}

//MemQuery Extract memory information according to StatRequest
func (s *Stat) MemQuery(ctx context.Context, req *StatRequest) (*MemReply, error) {

	idFieldName, nameFieldName := getIDNameFields(req)
	query := buildQueryString(req, idFieldName, "mem")
	res, err := s.conn.Query(query)
	if err != nil {
		return nil, err
	}

	memReply := MemReply{}
	list := res.Results[0].Series[0].Values
	memReply.Entries = make([]*MemEntry, len(list))
	for i, row := range list {
		entry, err := s.ressourceMemQuery(ctx, req, idFieldName, nameFieldName, row[1].(string))
		if err != nil {
			return nil, err
		}
		memReply.Entries[i] = entry
	}
	return &memReply, nil
}

//Extract memory information for one ressource according to StatRequest
func (s *Stat) ressourceMemQuery(ctx context.Context, req *StatRequest, idFieldName string, nameFieldName, ID string) (*MemEntry, error) {
	var memfields = "total_cache, rss, usage FROM docker_container_mem"
	ql := fmt.Sprintf("SELECT %s, %s WHERE %s='%s' ORDER BY time DESC LIMIT 1", nameFieldName, memfields, idFieldName, ID)
	res, err := s.conn.Query(ql)
	if err != nil {
		return nil, err
	}
	data := res.Results[0].Series[0].Values[0]
	entry := MemEntry{
		Id:     ID,
		Name:   data[0].(string),
		Cache: 	data[1].(int64),
		Rss:   	data[2].(int64),
		Usage: 	data[3].(int64),
	}
	return &entry, nil
}

//IOQuery Extract IO information according to CPURequest
func (s *Stat) IOQuery(ctx context.Context, req *StatRequest) (*IOReply, error) {

	idFieldName, nameFieldName := getIDNameFields(req)
	query := buildQueryString(req, idFieldName, "blkio")
	res, err := s.conn.Query(query)
	if err != nil {
		return nil, err
	}

	ioReply := IOReply{}
	list := res.Results[0].Series[0].Values
	ioReply.Entries = make([]*IOEntry, len(list))
	for i, row := range list {
		entry, err := s.ressourceIOQuery(ctx, req, idFieldName, nameFieldName, row[1].(string))
		if err != nil {
			return nil, err
		}
		ioReply.Entries[i] = entry
	}
	return &ioReply, nil
}

//Extract IO information for one ressource according to StatRequest
func (s *Stat) ressourceIOQuery(ctx context.Context, req *StatRequest, idFieldName string, nameFieldName, ID string) (*IOEntry, error) {
	var ioFields = "io_serviced_recursive_read, io_serviced_recursive_write, io_serviced_recursive_total, io_service_bytes_recursive_read, io_service_bytes_recursive_write, io_service_bytes_recursive_total"
	ql := fmt.Sprintf("SELECT %s, %s WHERE %s='%s' ORDER BY time DESC LIMIT 1", nameFieldName, ioFields, idFieldName, ID)
	res, err := s.conn.Query(ql)
	if err != nil {
		return nil, err
	}
	data := res.Results[0].Series[0].Values[0]
	entry := IOEntry{
		Id:     			ID,
		Name:   			data[0].(string),
		NumberRead: 	data[1].(int64),
		NumberWrite:  data[2].(int64),
		NumberTotal: 	data[3].(int64),
		SizeRead:			data[4].(int64),
		SizeWrite:		data[5].(int64),
		SizeTotal:		data[6].(int64),
	}
	return &entry, nil
}

//NetQuery Extract net information according to StatRequest
func (s *Stat) NetQuery(ctx context.Context, req *StatRequest) (*NetReply, error) {
	idFieldName, nameFieldName := getIDNameFields(req)
	query := buildQueryString(req, idFieldName, "net")
	res, err := s.conn.Query(query)
	if err != nil {
		return nil, err
	}

	netReply := NetReply{}
	list := res.Results[0].Series[0].Values
	netReply.Entries = make([]*NetEntry, len(list))
	for i, row := range list {
		entry, err := s.ressourceNetQuery(ctx, req, idFieldName, nameFieldName, row[1].(string))
		if err != nil {
			return nil, err
		}
		netReply.Entries[i] = entry
	}
	return &netReply, nil
}

//Extract net information for one ressource according to StatRequest
func (s *Stat) ressourceNetQuery(ctx context.Context, req *StatRequest, idFieldName string, nameFieldName, ID string) (*NetEntry, error) {
	var netFields = "rx_bytes, rx_errors, tx_bytes, tx_erros"
	ql := fmt.Sprintf("SELECT %s, %s WHERE %s='%s' ORDER BY time DESC LIMIT 1", nameFieldName, netFields, idFieldName, ID)
	res, err := s.conn.Query(ql)
	if err != nil {
		return nil, err
	}
	data := res.Results[0].Series[0].Values[0]
	entry := NetEntry{
		Id:     		ID,
		Name:   		data[0].(string),
		RxBytes: 		data[1].(int64),
		RxErrors:  	data[2].(int64),
		TxBytes: 		data[3].(int64),
		TxErrors:		data[4].(int64),
	}
	return &entry, nil
}

//Return specific field name for influx query concidering StatRequest
func getIDNameFields(req *StatRequest) (string, string) {
	var idFieldName = "com.docker.swarm.node.id"
	var nameFieldName = "host"
	if req.Discriminator == "container" {
		idFieldName = "container_id"
		nameFieldName="container_name"
	} else if req.Discriminator == "service" {
		idFieldName = "com.docker.swarm.service.id"
		nameFieldName = "com.docker.swarm.service.name"
	} else {
		req.Discriminator = "node"
	}
	return idFieldName, nameFieldName
}


//Compute the influx 'sql' query string concidering StatRequest
func buildQueryString(req *StatRequest, ressourceFieldName string, metric string) string {
	query := fmt.Sprintf("SELECT distinct %s FROM docker_container_%s WHERE container_id!=''", ressourceFieldName, metric)
	if (req.FilterDatacenter != "") {
		query+= fmt.Sprintf(" AND datacenter='%s'", req.FilterDatacenter)
	}
	if (req.FilterHost != "") {
		query+= fmt.Sprintf(" AND host='%s'", req.FilterHost)
	}
	if (req.FilterContainerId != "") {
		query+= fmt.Sprintf(" AND container_id='%s'", req.FilterContainerId)
	}
	if (req.FilterContainerName != "") {
		query+= fmt.Sprintf(" AND container_name='%s'", req.FilterContainerName)
	}
	if (req.FilterContainerImage != "") {
		query+= fmt.Sprintf(" AND container_image='%s'", req.FilterContainerImage)
	}
	if (req.FilterServiceId != "") {
		query+= fmt.Sprintf(" AND com.docker.swarm.service.id='%s'", req.FilterServiceId)
	}	
	if (req.FilterServiceName != "") {
		query+= fmt.Sprintf(" AND com.docker.swarm.service.name='%s'", req.FilterServiceName)
	}
	if (req.FilterNodeId != "") {
		query+= fmt.Sprintf(" AND com.docker.swarm.node.id='%s'", req.FilterNodeId)
	}
	return query
}

