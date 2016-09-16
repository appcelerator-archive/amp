package service

// import (
// 	"testing"
// )
//
// var (
// 	goodTestPorts = map[string]PortSpec{
// 		":3000":         {TargetPort: 3000},
// 		"80:3000":       {TargetPort: 3000, PublishedPort: 80},
// 		"80:3000/tcp":   {TargetPort: 3000, PublishedPort: 80, Protocol: "tcp"},
// 		"4000:3000/udp": {TargetPort: 3000, PublishedPort: 4000, Protocol: "udp"},
// 	}
//
// 	badTestPorts = map[string]string{
// 		"3000":      "port requires leading colon",
// 		"80:":       "published port must specify target port",
// 		"80:3000/":  "missing protocol",
// 		"80:3000/x": "incorrect protocol",
// 	}
// )
//
// // The testing isn't meant to be exhaustive since bad ports will fail anyway,
// // the tests just verify that we can generally identify the PortSpec components
// func TestRegex(t *testing.T) {
// 	for spec, expected := range goodTestPorts {
// 		//fmt.Println(spec)
// 		portSpec, err := ParsePortSpec(spec)
// 		if err != nil {
// 			t.Errorf("failed to parse \"%s\"\n%v", spec, err)
// 		}
//
// 		if portSpec.Name != expected.Name {
// 			t.Errorf("expected Name=%s, got: %s", expected.Name, portSpec.Name)
// 		}
//
// 		if portSpec.TargetPort != expected.TargetPort {
// 			t.Errorf("expected TargetPort=%d, got: %d", expected.TargetPort, portSpec.TargetPort)
// 		}
//
// 		if portSpec.PublishedPort != expected.PublishedPort {
// 			t.Errorf("expected PublishedPort=%d, got: %d", expected.PublishedPort, portSpec.PublishedPort)
// 		}
//
// 		if portSpec.Protocol != expected.Protocol {
// 			t.Errorf("expected Protocol=%s, got: %s", expected.Protocol, portSpec.Protocol)
// 		}
// 	}
//
// 	for spec, reason := range badTestPorts {
// 		//fmt.Printf("%s => %s\n", spec, reason)
// 		_, err := ParsePortSpec(spec)
// 		if err == nil {
// 			t.Errorf("should have failed to parse \"%s\" because %s", spec, reason)
// 		}
// 	}
// }
