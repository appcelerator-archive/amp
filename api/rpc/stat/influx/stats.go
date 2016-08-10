package influx

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

// Stats is the wrapper for the influx client connection
type Stats struct {
	client client.Client
	dbname string
	conn   string
	//TODO figure out better Security
	u string
	p string
}

// New return a newly created stats structure
func New(connection, dbname, u, p string) Stats {
	return Stats{
		dbname: dbname,
		conn:   connection,
		u:      u,
		p:      p,
	}
}

func (s *Stats) query(query string, database string) client.Query {
	return client.Query{
		Command:   query,
		Database:  database,
		Precision: "s",
	}
}

// Query executes the provided query string and returns the results as a JSON object
func (s *Stats) Query(q string) (string, error) {
	// ExecuteQuery runs any query statement
	response, err := s.client.Query(s.query(q, s.dbname))
	if err != nil {
		fmt.Printf("ERR: %s\n", err)
		return "", err
	}
	if err = response.Error(); err != nil {
		fmt.Printf("ERR: %s\n", response.Error())
		return "", err
	}
	data, err := json.Marshal(response)
	return string(data), err
}

// Connect to stats server
func (s *Stats) Connect(timeout time.Duration) error {
	// Make client
	//TODO Security!
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     s.conn,
		Username: s.u,
		Password: s.p,
	})

	if err != nil {
		return err
	}
	s.client = c
	return err
}

// Close connection to stats server
func (s *Stats) Close() error {
	err := s.client.Close()
	return err
}
