package stack

import (
	"io/ioutil"
	"path"
	"reflect"
	"strings"
	"testing"
	//	. "github.com/appcelerator/amp/api/rpc/stack"
	//	"context"
)

type TestSpec struct {
	fileName string
	valid    bool
	contents []byte
}

var (
	testDir = "./test_samples"

	sample1 = map[string]serviceMap{
		"pinger": {
			Image:    "appcelerator/pinger",
			Replicas: 2,
			Public: []publishSpec{
				{
					Name:         "www",
					PublishPort:  90,
					InternalPort: 3000,
					Protocol:     "tcp",
				},
			},
		},
	}

	sample2 = map[string]serviceMap{
		"web": {
			Image:    "appcelerator/amp-demo-service",
			Replicas: 3,
			Public: []publishSpec{
				{
					Name:         "www",
					PublishPort:  90,
					InternalPort: 3000,
					Protocol:     "tcp",
				},
			},
			Environment: map[string]string{
				"REDIS_PASSWORD": "password",
			},
		},
		"redis": {
			Image: "redis",
			Environment: map[string]string{
				"PASSWORD": "password",
			},
		},
	}
	sample3 = map[string]serviceMap{
		"pinger": {
			Image:    "appcelerator/pinger",
			Replicas: 2,
		},
		"pinger2": {
			Image:    "appcelerator/pinger",
			Replicas: 2,
			Public: []publishSpec{
				{
					Name:         "www",
					InternalPort: 3000,
					Protocol:     "tcp",
				},
			},
		},
	}
	sample4 = map[string]serviceMap{
		"python": {
			Image:    "tutum/quickstart-python",
			Replicas: 3,
			Public: []publishSpec{
				{
					Name:         "python",
					InternalPort: 80,
				},
			},
		},
		"go": {
			Image:    "htilford/go-redis-counter",
			Replicas: 3,
			Public: []publishSpec{
				{
					Name:         "go",
					InternalPort: 80,
				},
			},
		},
		"redis": {
			Image: "redis",
		},
	}

	// map of filenames to a map of serviceMap elements (each file has one or more)
	compareStructs = map[string]map[string]serviceMap{
		"sample-01.yml":  sample1,
		"sample-02.yml":  sample2,
		"sample-03.yml":  sample3,
		"sample-03.json": sample3,
		"sample-04.yml":  sample4,
	}
)

func TestSamples(t *testing.T) {
	tests := loadFiles(t)
	for _, test := range tests {
		parse(t, test)
	}
}

func loadFiles(t *testing.T) []*TestSpec {
	tests := []*TestSpec{}
	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		t.Error(err)
		return nil
	}

	for _, f := range files {
		name := f.Name()
		valid := false
		if !strings.HasPrefix(name, "invalid-") {
			valid = true
		}
		contents, err := ioutil.ReadFile(path.Join(testDir, name))
		if err != nil {
			t.Errorf("unable to load test sample: %s. Error: %v", name, err)
		}
		testSpec := &TestSpec{
			fileName: name,
			contents: contents,
			valid:    valid,
		}
		tests = append(tests, testSpec)
	}

	return tests
}

func parse(t *testing.T, test *TestSpec) {
	serviceSpecMap, err := parseAsServiceMap(test.contents)
	if err != nil {
		t.Error(err)
		return
	}

	for name, spec := range serviceSpecMap {
		if !spec.compare(t, compareStructs[test.fileName][name]) {
			t.Logf("name: %s, valid: %t, contents:\n%s", test.fileName, test.valid, string(test.contents))
			t.Log(serviceSpecMap)
			t.Errorf("FAIL: %s (%s)", name, test.fileName)
		}
	}

	// out, err := NewStackfromYaml(context.Background(), test.contents)
	// if err != nil {
	// 	t.Logf("fatal error parsing contents of %s: %v", test.fileName, err)
	// }
	// t.Logf("parsed => \n%v", out)
}

func (a serviceMap) compare(t *testing.T, b serviceMap) bool {
	if a.Image != b.Image {
		t.Logf("Images don't match: %v != %v\n", a.Image, b.Image)
		return false
	}
	if a.Replicas != b.Replicas {
		t.Logf("Replicas don't match: %v != %v\n", a.Replicas, b.Replicas)
		return false
	}
	if len(a.Public) != len(b.Public) {
		t.Logf("Public don't match: %v != %v\n", a.Public, b.Public)
		return false
	}
	for _, publishSpec := range a.Public {
		if !contains(b.Public, publishSpec) {
			t.Logf("Public doesn' contains: %v => %v\n", a.Public, b.Public)
			return false
		}
	}
	if !compareEnvironment(a, b) {
		t.Logf("Env don't match: %v != %v\n", a, b)
		return false
	}
	return true
}

// contains checks to see if a publishSpec is contained in a publishSpec slice
func contains(specs []publishSpec, spec publishSpec) bool {
	for _, cmp := range specs {
		if spec.compare(cmp) {
			return true
		}
	}
	return false
}

// compare returns true if and only if the members of both publishSpecs are equal
func (a publishSpec) compare(b publishSpec) bool {
	if a.Name != b.Name {
		return false
	}
	if a.PublishPort != b.PublishPort {
		return false
	}
	if a.InternalPort != b.InternalPort {
		return false
	}
	if a.Protocol != b.Protocol {
		return false
	}
	return true
}

// compareEnvironment returns true if and only if both serviceMap maps are equal (performs deep equal check)
func compareEnvironment(a serviceMap, b serviceMap) bool {
	ae := environmentToMap(a.Environment)
	be := environmentToMap(b.Environment)
	return reflect.DeepEqual(ae, be)
}

// environmentToMap
func environmentToMap(env interface{}) map[string]string {
	es, ok := env.(map[string]string)
	if ok {
		return es
	}

	envmap := make(map[string]string)

	em, ok := env.(map[interface{}]interface{})
	if ok {
		for k, v := range em {
			envmap[k.(string)] = v.(string)
		}
	}
	ea, ok := env.([]interface{})
	if ok {
		for _, s := range ea {
			a := strings.Split(s.(string), "=")
			k := a[0]
			v := a[1]
			envmap[k] = v
		}
	}

	return envmap
}
