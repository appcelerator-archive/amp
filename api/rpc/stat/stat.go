package stat

import (
	"time"

	"github.com/appcelerator/amp/data/influx"
	"golang.org/x/net/context"
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

//New retrun a new implementation of StateServer
func New(cfg Config) (*Stat, error) {

	c := influx.New(cfg.Connstr, cfg.Dbname, cfg.U, cfg.P)
	err := c.Connect(5 * time.Second)
	return &Stat{conn: c}, err
}

// ExecuteQuery implements business logic for StatServer interface
func (s *Stat) ExecuteQuery(ctx context.Context, req *QueryRequest) (*QueryReply, error) {
	resp, err := s.conn.Query(req.Query)
	if err != nil {
		return nil, err
	}
	data, err := s.conn.Marshal(resp)
	if err != nil {
		return nil, err
	}
	rep := &QueryReply{Response: data}
	return rep, err
}
