package stat

import (
	"github.com/influxdata/influxdb/client/v2"
)

type Marshaler interface {
	setFieldValue(string, string)
	nextSlice()
}

type CPUContainerMarshaler struct {
	cpus   []*CPUByContainer
	offSet int
}

func NewCPUContainerMarshaler(size int) *CPUContainerMarshaler {
	marshaler := CPUContainerMarshaler{}

	marshaler.cpus = make([]*CPUByContainer, size)
	marshaler.offSet = 0
	marshaler.cpus[0] = &CPUByContainer{}
	return &marshaler
}
func (c *CPUContainerMarshaler) nextSlice() {
	c.offSet++
	if c.offSet < len(c.cpus) {
		c.cpus[c.offSet] = &CPUByContainer{}
	}
}

func (c *CPUContainerMarshaler) setFieldValue(f string, v string) {

	switch f {
	case "container_image":
		c.cpus[c.offSet].Containerimage = v
	case "container_name":
		c.cpus[c.offSet].Containername = v
	case "datacenter":
		c.cpus[c.offSet].Datacenter = v
	case "time":
		c.cpus[c.offSet].Time = v
	case "cpu_pct":
		c.cpus[c.offSet].Cpupct = v
	default:
		//ignore
	}
}

// MarshalInfluxToProto takes the influx result structure and converts it to protobut message
func MarshalInfluxToProto(resp *client.Response, m Marshaler) {

	for i := 0; len(resp.Results) > 0 && i < len(resp.Results[0].Series); i++ {

		for k, v := range resp.Results[0].Series[i].Tags {
			m.setFieldValue(k, v)
		}
		for j := 0; j < len(resp.Results[0].Series[i].Columns); j++ {
			//TODO Need to add type support
			m.setFieldValue(resp.Results[0].Series[i].Columns[j], "TODO")
		}
		m.nextSlice()
	}
}
