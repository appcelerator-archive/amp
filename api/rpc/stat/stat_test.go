package stat

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/appcelerator/amp/data/influx"

	"golang.org/x/net/context"
)

const (
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

	err := srv.Influx.Connect(5 * time.Second)
	if err != nil {
		panic(err)
	}
	defer srv.Influx.Close()
	os.Exit(m.Run())
}

// Excercises the rpc service for CPUQuery
func TestCPUQueryServiceWithLimit(t *testing.T) {
	//Build a query to validate the following command from the CLI
	//amp stats --cpu --container --service-name=kafka --period 5m
	query := StatRequest{}
	//set discriminator
	query.Discriminator = "container"
	//Set filters
	query.FilterServiceName = "kafka"
	query.Period = "5m"
	query.Limit = "1"
	// Call the service over the call stack directly
	res, err := srv.CPUQuery(context.Background(), &query)
	if err != nil {
		t.Error(err)
	}
	if len(res.Entries) != 1 {
		t.Errorf("Unexpected a response from influx %v\n", len(res.Entries))
	}

}

// Excercises the rpc service for CPUQuery
func TestCPUQueryService(t *testing.T) {
	//Build a query to validate the following command from the CLI
	//amp stats --cpu --container --service-name=kafka --period 5m
	query := StatRequest{}
	//set discriminator
	query.Discriminator = "container"
	//Set filters
	query.FilterServiceName = "kafka"
	query.Period = "5m"
	// Call the service over the call stack directly
	res, err := srv.CPUQuery(context.Background(), &query)
	if err != nil {
		t.Error(err)
	}
	if len(res.Entries) == 0 {
		t.Errorf("Expected a response from influx \n")
	}

}
func createStatServer() *Stat {
	//Create the config
	var stat = &Stat{}
	host := os.Getenv("influxhost")
	cstr := "http://localhost:8086"
	if host != "" {
		cstr = "http://" + host + ":8086"
	}
	stat.Influx = influx.New(cstr, "telegraf", "", "")
	return stat
}
