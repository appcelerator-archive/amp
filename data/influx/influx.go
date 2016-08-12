package influx

import (
	"encoding/json"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

// Influx is the wrapper for the influx client connection
type Influx struct {
	client client.Client
	//Influx can support multiple databases within a single cluster
	//Each connection should specify which database to work with eg, "telegraf"
	dbname string
	conn   string
	//TODO figure out better Security
	u string
	p string
}

// New returns a newly created influxdb object
func New(connection, dbname, u, p string) Influx {
	return Influx{
		dbname: dbname,
		conn:   connection,
		u:      u,
		p:      p,
	}
}

func (s *Influx) query(query string, database string) client.Query {
	return client.Query{
		Command:   query,
		Database:  database,
		Precision: "s",
	}
}

// Query executes the provided query string and returns the results as a JSON object
func (s *Influx) Query(q string) (string, error) {
	// ExecuteQuery runs any query statement
	response, err := s.client.Query(s.query(q, s.dbname))
	if err != nil {
		return "", err
	}
	if err = response.Error(); err != nil {
		return "", err
	}
	data, err := json.Marshal(response)
	return string(data), err
}

// Connect to influxdb server
func (s *Influx) Connect(timeout time.Duration) error {
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

// Close connection to influxdb server
func (s *Influx) Close() error {
	err := s.client.Close()
	return err
}
