package influx

import (
	"os"
	"strings"
	"testing"
)

var (
	s    Stats
	host string
)

func TestMain(m *testing.M) {
	statsInit()
	defer s.Close()
	os.Exit(m.Run())
}

func TestQuery(t *testing.T) {
	res, err := s.Query("SHOW MEASUREMENTS")
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(res, "measurements") {
		t.Errorf("Expected String to contain %s, actual=%s \n", "Measurement", res)
	}

}

func statsInit() {
	host := os.Getenv("influxhost")
	cstr := "http://localhost:8086"
	if host != "" {
		cstr = "http://" + host + ":8086"
	}
	s = New(cstr, "_internal", "admin", "changme")
	s.Connect(5)

}
