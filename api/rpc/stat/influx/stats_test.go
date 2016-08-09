package influx

import (
	"os"
	"strings"
	"testing"

	"github.com/appcelerator/amp/api/rpc/stat"
)

var (
	s stat.Stats
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
	s = New("http://localhost:8086", "_internal", "admin", "changme")
	s.Connect(5)
}
