package cli_test

import (
	"github.com/appcelerator/amp/api/server"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"
)

type TestSpec struct {
	fileName string
	contents []byte
}

type commandSpec struct {
	Cmd     string   `yaml:"cmd"`
	Args    []string `yaml:"args"`
	Options []string `yaml:"options"`
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
		t.Log("Loading file:", name)
		contents, err := ioutil.ReadFile(path.Join(testDir, name))
		if err != nil {
			t.Errorf("unable to load test sample: %s. Error: %v", name, err)
		}
		testSpec := &TestSpec{
			fileName: name,
			contents: contents,
		}
		tests = append(tests, testSpec)
	}
	return tests
}

func parse(t *testing.T, test *TestSpec) {
	commandMap, err := parseCommandMap(test.contents)
	if err != nil {
		t.Error(err)
		return
	}
	for _, spec := range commandMap {
		cmdSplit := strings.Fields(spec.Cmd)
		optionsSplit := []string{}
		for _, val := range spec.Options {
			optionsSplit = append(optionsSplit, strings.Fields(val)...)
		}
		commandFinal := append(cmdSplit, spec.Args...)
		commandFinal = append(commandFinal, optionsSplit...)
		t.Log(commandFinal, "Command passed.")
		runCmd(t, commandFinal)
	}
}

func runCmd(t *testing.T, cmdString []string) {
	t.Log(cmdString, "Running...")
	for i := 0; i < 10; i++ {
		t.Log(cmdString, " Iteration:", i+1)
		cmd := exec.Command(cmdString[0], cmdString[1:]...)
		result, err := cmd.CombinedOutput()
		if err != nil && i >= 9 {
			t.Log(cmdString, " Error:", err)
			t.Log(cmdString, " Command has failed, exiting.")
			t.Fail()
		} else if err != nil {
			t.Log(cmdString, " Error:", err)
			t.Log(cmdString, " Command failed, retrying.")
			time.Sleep(1 * time.Second)
		} else {
			t.Log(cmdString, " Command result:\n", string(result))
			break
		}
	}
}

func parseCommandMap(b []byte) (out map[string]commandSpec, err error) {
	err = yaml.Unmarshal(b, &out)
	return
}
