package influx

import (
	"os"
	"testing"
	"time"
)

var (
	influx Influx
)

func TestMain(m *testing.M) {
	influxInit()
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

func influxInit() {
	host := os.Getenv("influxhost")
	cstr := "http://localhost:8086"
	if host != "" {
		cstr = "http://" + host + ":8086"
	}
	influx = New(cstr, "_internal", "admin", "changme")
	influx.Connect(time.Second * 60)
}
