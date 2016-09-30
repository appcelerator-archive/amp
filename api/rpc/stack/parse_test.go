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

	sample1 = map[string]serviceSpec{
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

	sample2 = map[string]serviceSpec{
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
	sample3 = map[string]serviceSpec{
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
		"haproxy": {
  			Public: []publishSpec{
  				{
    					PublishPort: 83,
      					InternalPort: 80,

  				},
  			},
  		},
	}
	sample4 = map[string]serviceSpec{
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

	sample5 = map[string]serviceSpec{
		"pinger": {
			Image: "appcelerator/pinger",
			Environment: map[string]string{
				"foo": "bar",
			},
			Public: []publishSpec{
				{
					PublishPort:  3000,
					InternalPort: 3000,
				},
			},
		},
	}

	sample6 = map[string]serviceSpec{
		"pinger": {
			Image: "appcelerator/pinger",
			Labels: map[string]string{
				"foo":   "bar",
				"hello": "world",
			},
			Public: []publishSpec{
				{
					PublishPort:  3000,
					InternalPort: 3000,
				},
			},
		},
	}

	sample7 = map[string]serviceSpec{
		"pinger": {
			Image: "appcelerator/pinger",
			ContainerLabels: map[string]string{
				"foo":   "bar",
				"hello": "world",
			},
			Public: []publishSpec{
				{
					PublishPort:  3000,
					InternalPort: 3000,
				},
			},
		},
	}

	sample8_1 = map[string]serviceSpec{
		"pinger": {
			Image:    "appcelerator/pinger",
			Mode:     "replicated",
			Replicas: 3,
			Public: []publishSpec{
				{
					PublishPort:  3000,
					InternalPort: 3000,
				},
			},
		},
	}

	sample8_2 = map[string]serviceSpec{
		"pinger": {
			Image: "appcelerator/pinger",
			Mode:  "global",
			Public: []publishSpec{
				{
					PublishPort:  3000,
					InternalPort: 3000,
				},
			},
		},
	}

	// map of filenames to a map of serviceSpec elements (each file has one or more)
	compareSpecs = map[string]map[string]serviceSpec{
		"sample-01.yml":                    sample1,
		"sample-02.yml":                    sample2,
		"sample-03.yml":                    sample3,
		"sample-03.json":                   sample3,
		"sample-04.yml":                    sample4,
		"sample-05-1-env.yml":              sample5,
		"sample-05-2-env.yml":              sample5,
		"sample-06-1-service-labels.yml":   sample6,
		"sample-06-2-service-labels.yml":   sample6,
		"sample-07-1-container-labels.yml": sample7,
		"sample-07-2-container-labels.yml": sample7,
		"sample-08-1-mode.yml":             sample8_1,
		"sample-08-2-mode.yml":             sample8_2,
	}
)

func TestSamples(t *testing.T) {
	tests := loadFiles(t)
	for _, test := range tests {
		if compareSpecs[test.fileName] == nil {
			t.Logf("WARNING: skipping '%s' because the comparison sample is missing", test.fileName)
			continue
		}
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
	serviceMap, err := parseServiceMap(test.contents)
	if err != nil {
		t.Error(err)
		return
	}

	for name, spec := range serviceMap {
		if !spec.compare(t, compareSpecs[test.fileName][name]) {
			t.Logf("name: %s, valid: %t, contents:\n%s", test.fileName, test.valid, string(test.contents))
			t.Log(serviceMap)
			t.Errorf("FAIL: %s (%s)", name, test.fileName)
		}
	}
}

// compares actual (parsed) service spec to expected
func (a serviceSpec) compare(t *testing.T, b serviceSpec) bool {
	if a.Image != b.Image {
		t.Logf("actual != expected (image): '%v' != '%v'\n", a.Image, b.Image)
		return false
	}
	if a.Replicas != b.Replicas {
		t.Logf("actual != expected (replicas): %v != %v\n", a.Replicas, b.Replicas)
		return false
	}
	if !reflect.DeepEqual(a.Public, b.Public) {
		t.Logf("actual != expected (public): %v != %v\n", a.Public, b.Public)
		return false
	}
	if !reflect.DeepEqual(toMap(a.Environment), toMap(b.Environment)) {
		t.Logf("actual != expected (env): %v != %v\n", a.Environment, b.Environment)
		return false
	}
	if !reflect.DeepEqual(toMap(a.Labels), toMap(b.Labels)) {
		t.Logf("actual != expected (labels): %v != %v\n", a.Labels, b.Labels)
		return false
	}
	if !reflect.DeepEqual(toMap(a.ContainerLabels), toMap(b.ContainerLabels)) {
		t.Logf("actual != expected (container_labels): %v != %v\n", a.ContainerLabels, b.ContainerLabels)
		return false
	}
	return true
}

// compareEnvironment returns true if and only if both serviceSpec maps are equal (performs deep equal check)
// first, ensure both environments are converted to maps (because an environment might be a string array or a map)
func compareEnvironment(a serviceSpec, b serviceSpec) bool {
	ae := toMap(a.Environment)
	be := toMap(b.Environment)
	return reflect.DeepEqual(ae, be)
}

// toMap attempts to cast an interface to either an array of strings or a map of strings to strings and then returns a typed map
func toMap(arrayOrMap interface{}) map[string]string {
	stringmap, ok := arrayOrMap.(map[string]string)
	if ok {
		return stringmap
	}

	stringmap = make(map[string]string)

	if m, ok := arrayOrMap.(map[interface{}]interface{}); ok {
		for k, v := range m {
			stringmap[k.(string)] = v.(string)
		}
	} else if a, ok := arrayOrMap.([]interface{}); ok {
		for _, s := range a {
			parts := strings.Split(s.(string), "=")
			stringmap[parts[0]] = parts[1]
		}
	}

	return stringmap
}
