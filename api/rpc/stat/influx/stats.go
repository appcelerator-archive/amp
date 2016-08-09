package influx

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

// New return a newly created stats structure
func New(connection, dbname string) *stats {
	return &stats{
		dbname: dbname,
		conn:   connection,
	}
}

type stats struct {
	client client.Client
	dbname string
	conn   string
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
	s.writeJSON(response, os.Stdout)
	return "", err
}
func (s *stats) Endpoints() []string {
	return nil
}

// Connect to stats server
func (s *stats) Connect(timeout time.Duration) (*stats, error) {
	// Make client
	//TODO Security!
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     s.conn,
		Username: "admin",
		Password: "changeme",
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return nil, err
	}
	s.client = c
	return s, err
}

// Close connection to stats server
func (s *stats) Close() error {
	err := s.client.Close()
	if err != nil {
		return err
	}
	return err
}

func (s *stats) writeJSON(response *client.Response, w io.Writer) {
	var data []byte
	var err error

	data, err = json.MarshalIndent(response, "", "    ")
	if err != nil {
		fmt.Fprintf(w, "Unable to parse json: %s\n", err)
		return
	}
	fmt.Fprintln(w, string(data))
}
