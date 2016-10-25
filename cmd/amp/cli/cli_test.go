package cli_test

import (
	"github.com/appcelerator/amp/api/server"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"testing"
	"time"
)

type TestSpec struct {
	fileName string
	contents []byte
	valid    bool
}

type CommandSpec struct {
	Cmd         string   `yaml:"cmd"`
	Args        []string `yaml:"args"`
	Options     []string `yaml:"options"`
	Expectation string   `yaml:"expectation"`
}

var (
	testDir = "./test_samples"
)

func TestMain(t *testing.T) {
	_, conn := server.StartTestServer()
	t.Log(conn)
}

func TestCmds(t *testing.T) {
	tests := loadFiles(t)
	for _, test := range tests {
		t.Log("-----------------------------------------------------------------------------------------")
		t.Logf("test %s\n", test.fileName)
		parseCmd(t, test)
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
		t.Log("Loading file:", name)
		valid := false
		if !strings.HasPrefix(name, "00-") {
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

func parseCmd(t *testing.T, test *TestSpec) {
	commandMap, err := generateCmdSpec(test.contents)
	if err != nil {
		t.Error(err)
		return
	}
	for _, cmdSpec := range commandMap {
		cmdString := generateCmdString(cmdSpec)
		t.Log(cmdString, "Command passed.")
		for i := 0; i < 10; i++ {
			t.Log(cmdString, "Running...")
			t.Log(cmdString, "Iteration:", i+1)
			result, err := runCmd(cmdString)
			validID := regexp.MustCompile(cmdSpec.Expectation)
			if test.valid == false {
				if err == nil {
					t.Log(cmdString, "Error:", err)
					t.Log(cmdString, "Invalid Sample Command has failed, retrying.")
					time.Sleep(1 * time.Second)
				} else {
					if !validID.MatchString(string(result)) {
						t.Log(cmdString, "Error: miss matched expectation")
						t.Fail()
						break
					}
					t.Log(cmdString, "Invalid Sample Command result:\n", string(result))
					break
				}
			} else {
				if err != nil {
					t.Log(cmdString, "Error:", err)
					t.Log(cmdString, "Command failed, retrying.")
					time.Sleep(1 * time.Second)
				} else {
					if !validID.MatchString(string(result)) {
						t.Log(cmdString, "Error: miss matched expectation")
						t.Fail()
						break
					}
					t.Log(cmdString, "Command result:\n", string(result))
					break
				}
			}
			if i >= 9 {
				t.Log(cmdString, "Error:", err)
				t.Log(cmdString, "Command has failed, exiting.")
				t.Fail()
			}
		}
	}
}

func generateCmdSpec(b []byte) (out map[string]CommandSpec, err error) {
	err = yaml.Unmarshal(b, &out)
	return
}

func generateCmdString(cmdSpec CommandSpec) (cmdString []string) {
	cmdSplit := strings.Fields(cmdSpec.Cmd)
	optionsSplit := []string{}
	for _, val := range cmdSpec.Options {
		optionsSplit = append(optionsSplit, strings.Fields(val)...)
	}
	cmdString = append(cmdSplit, cmdSpec.Args...)
	cmdString = append(cmdString, optionsSplit...)
	return
}

func runCmd(cmdString []string) (result []byte, err error) {
	cmd := exec.Command(cmdString[0], cmdString[1:]...)
	result, err = cmd.CombinedOutput()
	return
}
