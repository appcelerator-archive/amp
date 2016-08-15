package influx

import (
	"os"
	"strings"
	"testing"
)

var (
	influx Influx
	host   string
)

func TestMain(m *testing.M) {
	influxInit()
	defer influx.Close()
	os.Exit(m.Run())
}

func TestQuery(t *testing.T) {
	res, err := influx.Query("SHOW MEASUREMENTS")
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(res, "measurements") {
		t.Errorf("Expected String to contain %s, actual=%s \n", "Measurement", res)
	}

}

func influxInit() {
	host := os.Getenv("influxhost")
	cstr := "http://localhost:8086"
	if host != "" {
		cstr = "http://" + host + ":8086"
	}
	influx = New(cstr, "_internal", "admin", "changme")
	influx.Connect(5)

}
