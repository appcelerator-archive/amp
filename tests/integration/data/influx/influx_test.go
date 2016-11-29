package influx

import (
	"github.com/appcelerator/amp/config"
	. "github.com/appcelerator/amp/data/influx"
	"os"
	"testing"
	"time"
)

var (
	influx Influx
)

func TestMain(m *testing.M) {
	influx = New(amp.InfluxDefaultURL, "_internal", "admin", "changme")
	influx.Connect(60 * time.Second)
	defer influx.Close()
	os.Exit(m.Run())
}

func TestQuery(t *testing.T) {
	res, err := influx.Query("SHOW MEASUREMENTS")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Results) == 0 {
		t.Errorf("Expected results")
	}
	if len(res.Results[0].Series) == 0 {
		t.Errorf("Expected series")
	}
	if res.Results[0].Series[0].Name != "measurements" {
		t.Errorf("Expected name to be %s, actual=%s \n", "measurement", res.Results[0].Series[0].Name)
	}
}
