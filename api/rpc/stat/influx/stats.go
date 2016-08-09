package influx

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/appcelerator/amp/api/rpc/stat"
	"github.com/influxdata/influxdb/client/v2"
)

type stats struct {
	client client.Client
	dbname string
	conn   string
	//TODO figure out better Security
	u string
	p string
}

// New return a newly created stats structure
func New(connection, dbname, u, p string) stat.Stats {
	return &stats{
		dbname: dbname,
		conn:   connection,
		u:      u,
		p:      p,
	}
}

func (s *stats) query(query string, database string) client.Query {
	return client.Query{
		Command:   query,
		Database:  database,
		Precision: "s",
	}
}

// Query executes the provided query string and returns the results as a JSON object
func (s *stats) Query(q string) (string, error) {
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
func (s *stats) Endpoints() []string {
	return nil
}

// Connect to stats server
func (s *stats) Connect(timeout time.Duration) error {
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
func (s *stats) Close() error {
	err := s.client.Close()
	return err
}

// writeSJON takes the response and marshals it to JSON
func (s *stats) writeJSON(response *client.Response, w io.Writer) {
	var data []byte
	var err error
	data, err = json.Marshal(response)
	if err != nil {
		fmt.Fprintf(w, "Unable to parse json: %s\n", err)
		return
	}
	fmt.Fprintln(w, string(data))
}
