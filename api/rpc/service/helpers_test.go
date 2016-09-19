package service_test

import (
	"testing"
        . "github.com/appcelerator/amp/api/rpc/service"
)

var (
	goodPublishSpecs = []struct {
		input  string
		expect PublishSpec
	}{
		{input: ":3000", expect: PublishSpec{InternalPort: 3000}},
		{input: "80:3000", expect: PublishSpec{InternalPort: 3000, PublishPort: 80}},
		{input: "80:3000/tcp", expect: PublishSpec{InternalPort: 3000, PublishPort: 80, Protocol: "tcp"}},
		{input: "4000:3000/udp", expect: PublishSpec{InternalPort: 3000, PublishPort: 4000, Protocol: "udp"}},
		{input: "host:4000:3000/udp", expect: PublishSpec{InternalPort: 3000, PublishPort: 4000, Protocol: "udp"}},
		{input: "host-1:4000:3000/udp", expect: PublishSpec{InternalPort: 3000, PublishPort: 4000, Protocol: "udp"}},
	}

	badPublishSpecs = []struct {
		input  string
		reason string
	}{
		{input: "3000", reason: "internal port requires leading colon"},
		{input: "80:", reason: "publish port must specify internal port"},
		{input: "80:3000/", reason: "missing protocol"},
		{input: "80:3000/x", reason: "invalid protocol"},
		{input: "host-:4000:3000/udp", expect: PublishSpec{InternalPort: 3000, PublishPort: 4000, Protocol: "udp"}},
		{input: "-host:4000:3000/udp", expect: PublishSpec{InternalPort: 3000, PublishPort: 4000, Protocol: "udp"}},
	}
)

// This testing isn't meant to be rigorously exhaustive since bad ports will fail anyway,
// the tests just verify that we can generally identify PublishSpec components
func TestRegex(t *testing.T) {
	for _, test := range goodPublishSpecs {
                input := test.input
                expected := test.expect

		//fmt.Println(input)
		spec, err := ParsePublishSpec(input)
                t.Logf("input: %s, PublishSpec: %v", input, spec)

		if err != nil {
			t.Errorf("failed to parse \"%s\"\n%v", input, err)
		}

		if spec.Name != expected.Name {
			t.Errorf("expected Name=%s, got: %s", expected.Name, spec.Name)
		}

		if spec.InternalPort != expected.InternalPort {
			t.Errorf("expected InternalPort=%d, got: %d", expected.InternalPort, spec.InternalPort)
		}

		if spec.PublishPort != expected.PublishPort {
			t.Errorf("expected PublishPort=%d, got: %d", expected.PublishPort, spec.PublishPort)
		}

		if spec.Protocol != expected.Protocol {
			t.Errorf("expected Protocol=%s, got: %s", expected.Protocol, spec.Protocol)
		}
	}

	// TODO
	// for _, test := range badPublishSpecs {
        //         input := test.input
        //         reason := test.reason
        //
	// 	//fmt.Printf("%s => %s\n", input, reason)
	// 	_, err := ParsePublishSpec(input)
	// 	if err == nil {
	// 		t.Errorf("should have failed to parse \"%s\" because: %s", input, reason)
	// 	}
	// }
}
