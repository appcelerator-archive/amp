package stack

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strings"
	"testing"
	//	. "github.com/appcelerator/amp/api/rpc/stack"
	//	"context"
)

type TestSpec struct {
	fileName string
	valid    bool
	contents []byte
	ref      *stackSpec
}

var (
	testDir = "./test_samples"

	sample1 = stackSpec{
		Services: map[string]serviceSpec{
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
		},
	}

	sample2 = stackSpec{
		Services: map[string]serviceSpec{
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
				Networks: map[string]networkAliases{
					"app-net": {
						Aliases: []string{"stack1-redis"},
					},
				},
			},
		},
		Networks: map[string]networkSpec{
			"app-net": {
				Driver: "overlay",
			},
		},
	}

	sample3 = stackSpec{
		Services: map[string]serviceSpec{
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
						PublishPort:  83,
						InternalPort: 80,
					},
				},
			},
		},
	}

	sample4 = stackSpec{
		Services: map[string]serviceSpec{
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
		},
	}

	sample5 = stackSpec{
		Services: map[string]serviceSpec{
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
		},
	}

	sample6 = stackSpec{
		Services: map[string]serviceSpec{
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
		},
	}

	sample7 = stackSpec{
		Services: map[string]serviceSpec{
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
		},
	}

	sample8_1 = stackSpec{
		Services: map[string]serviceSpec{
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
		},
	}

	sample8_2 = stackSpec{
		Services: map[string]serviceSpec{
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
		},
	}

	// map of filenames to a map of serviceSpec elements (each file has one or more)
	compareSpecs = map[string]stackSpec{
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

// test all test samples
func TestSamples(t *testing.T) {
	tests := loadFiles(t)
	for _, test := range tests {
		t.Log("-----------------------------------------------------------------------------------------")
		t.Logf("test %s\n", test.fileName)
		parse(t, test)
	}
}

// process one test and verify if ok or not
func parse(t *testing.T, test *TestSpec) {
	parsedStack, err := parseStack(test.contents)
	//t.Logf("%+v\n", parsedStack)
	//t.Logf("%+v\n", test.ref)
	if err != nil {
		t.Error(err)
		return
	}
	diff, ok := parsedStack.extractDiff(t, *test.ref)
	//diff, ok := test.ref.extractDiff(t, *parsedStack)
	if !ok {
		t.Errorf("FAIL on service: %s\n", test.fileName)
		if diff[1] == "" {
			t.Logf("not expected, but in file: %+v\n", diff[0])
		} else if diff[0] == "" {
			t.Logf("expected, but not in file: %+v\n", diff[1])
		} else {
			t.Logf("expected:       %+v\n", diff[1])
			t.Logf("found in file:  %+v\n", diff[0])
		}
		//t.Logf("from file: %+v\n", parsedStack.Services[name])
	} else {
		t.Log("Tested ok")
	}
}

// return the most explicite two strings explaining the difference between the two structs
func (a stackSpec) extractDiff(t *testing.T, b stackSpec) ([2]string, bool) {
	//t.Log("process file")
	sa := fmt.Sprintf("%+v", a)
	la := explodeExtend(t, simplifyString(sa))
	//t.Log("process ref")
	sb := fmt.Sprintf("%+v", b)
	lb := explodeExtend(t, simplifyString(sb))
	if len(la) != len(lb) {
		return [2]string{
			getDiff(la, lb),
			getDiff(lb, la),
		}, false
	}
	for i, item := range la {
		if item != lb[i] {
			return [2]string{item, lb[i]}, false
		}
	}
	return [2]string{}, true
}

// supress not useful charactere and normalize map syntax
func simplifyString(line string) string {
	line = line[1 : len(line)-1]
	line = strings.Replace(line, "map[", "[", -1)
	line = strings.Replace(line, "{", "[", -1)
	line = strings.Replace(line, "}", "]", -1)
	line = strings.Replace(line, "=", ":", -1)
	line = strings.Replace(line, "&", "", -1)
	line = strings.Replace(line, "[]", "<nil>", -1)
	return line
}

// create a list of orderer struct member fullname:value to be able to be compare item by item
func explodeExtend(t *testing.T, line string) []string {
	//t.Logf("initial: %s\n", line)
	found := true
	for found {
		line, found = extendOrderNames(t, line)
		if found {
			//t.Logf("extended: %s\n", line)
		}
	}
	list := strings.Split(line, " ")
	sort.Strings(list)
	//t.Logf("result: %v\n", list)
	return list
}

//Give at each struc member name its full name with all the parent strucs name
func extendOrderNames(t *testing.T, line string) (string, bool) {
	//t.Log(line)
	i3 := strings.Index(line, "]")
	if i3 < 0 {
		return line, false
	}
	i2 := strings.LastIndex(line[0:i3], "[")
	if i2 < 0 {
		return line, false
	}
	if i2 > 0 && line[i2-1] != ':' {
		//t.Logf("i2=%d i3=%d no name", i2, i3)
		subline := line[i2+1 : i3]
		//t.Log(subline)
		line = line[0:i2] + subline + line[i3+1:]
		return line, true
	}
	name := ""
	i1 := strings.LastIndex(line[0:i2], " ")
	i1b := strings.LastIndex(line[0:i2], "[")
	if i1b > i1 {
		i1 = i1b
	}
	if i1 < 0 {
		//i1 = 0
		name = line[0:i2]
	} else {
		name = line[i1+1 : i2]
	}
	//t.Logf("i1=%d i2=%d i3=%d name=%s\n", i1, i2, i3, name)
	subline := line[i2+1 : i3]
	//t.Log(subline)
	list := strings.Split(subline, " ")
	sort.Strings(list)
	newSubline := ""
	for _, item := range list {
		if strings.HasPrefix(item, "{") {
			item = item[1:]
		}
		if strings.HasSuffix(item, "}") {
			item = item[0 : len(item)-2]
		}
		newSubline += (name + item + " ")
	}
	line = line[0:i1+1] + newSubline[0:len(newSubline)-1] + line[i3+1:]
	//t.Log(line)
	return line, true
}

// extract the list if the fullname:value of l1 which don't exist in l2, return catenate string for display
func getDiff(l1 []string, l2 []string) string {
	ret := ""
	for _, item1 := range l1 {
		found := false
		for _, item2 := range l2 {
			if item1 == item2 {
				found = true
			}
		}
		if !found {
			ret += item1 + " "
		}
	}
	return ret
}

// load files from samples dir
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
		spec, exist := compareSpecs[name]
		if !exist {
			t.Logf("WARNING: skipping '%s' because the comparison sample is missing", name)
		} else {
			testSpec := &TestSpec{
				fileName: name,
				contents: contents,
				valid:    valid,
				ref:      &spec,
			}
			tests = append(tests, testSpec)
		}
	}
	return tests
}
