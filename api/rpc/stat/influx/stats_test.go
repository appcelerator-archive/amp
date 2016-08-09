package influx

import "testing"

func TestQuery(t *testing.T) {
	s := New("http://localhost:8086", "_internal")
	s.Connect(5)
	s.Query("select * from runtime")
}
