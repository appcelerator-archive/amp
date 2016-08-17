package stat

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"
)

const (
	svcAddr  = "localhost:51001"
	queryVal = "measurements"
)

var (
	srv  *Stat
	host string
)

func TestMain(m *testing.M) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("test: ")
	// Create an instance of the service interface
	srv = createStatServer()
	err := srv.conn.Connect(5 * time.Second)
	if err != nil {
		panic(err)
	}
	defer srv.conn.Close()
	os.Exit(m.Run())
}

// Create a grpc client to invoke the stat service over specified port
// This runs in the same process but executes the call over the wire
func TestExecuteQuery(t *testing.T) {
	req := QueryRequest{Database: "", Query: "SHOW MEASUREMENTS"}

	// Call the service over the call stack directly
	res, err := srv.ExecuteQuery(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(res.Response, queryVal) {
		t.Errorf("Expected String to contain %s, actual=%s \n", queryVal, res.Response)
	}
}
func createStatServer() *Stat {
	//Create the config
	host := os.Getenv("influxhost")
	cstr := "http://localhost:8086"
	if host != "" {
		cstr = "http://" + host + ":8086"
	}

	cfg := Config{Connstr: cstr, Dbname: "_internal", U: "admin", P: "changeme"}
	// Create the service to allow service calls directly
	srv, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return srv
}
