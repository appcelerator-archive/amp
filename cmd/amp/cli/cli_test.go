package cli_test

import (
	"bytes"
	"fmt"
	"github.com/appcelerator/amp/api/server"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"text/template"
)

type TestSpec struct {
	Name     string
	Commands []CommandSpec
}

type CommandSpec struct {
	Cmd               string   `yaml:"cmd"`
	Args              []string `yaml:"args"`
	Options           []string `yaml:"options"`
	Expectation       string   `yaml:"expectation"`
	ExpectErrorStatus bool     `yaml:"expectErrorStatus"`
}

type LookupSpec struct {
	Name string
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	testDir   = "./test_samples"
	lookupDir = "./lookup"
	regexMap  map[string]string
)

func TestMain(m *testing.M) {
	server.StartTestServer()
	os.Exit(m.Run())
}

func TestCmds(t *testing.T) {
	err := loadRegexLookup()
	if err != nil {
		t.Errorf("Unable to load lookup specs, reason: %v", err)
		return
	}

	tests, err := loadTestSpecs()
	if err != nil {
		t.Errorf("unable to load test specs, reason: %v", err)
		return
	}

	for _, test := range tests {
		t.Log("-----------------------------------------------------------------------------------------")
		t.Logf("Running spec: %s", test.Name)
		if err := runTestSpec(t, test); err != nil {
			t.Error(err)
			return
		}
	}
}

func loadTestSpecs() ([]*TestSpec, error) {
	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		return nil, err
	}

	tests := []*TestSpec{}
	for _, file := range files {
		test, err := loadTestSpec(path.Join(testDir, file.Name()))
		if err != nil {
			return nil, err
		}
		if test != nil {
			tests = append(tests, test)
		}
	}
	return tests, nil
}

func loadTestSpec(fileName string) (*TestSpec, error) {
	if filepath.Ext(fileName) != ".yml" {
		return nil, nil
	}
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("unable to load test spec: %s. Error: %v", fileName, err)
	}
	testSpec := &TestSpec{
		Name: fileName,
	}

	var commandMap []CommandSpec
	if err := yaml.Unmarshal(content, &commandMap); err != nil {
		return nil, fmt.Errorf("unable to parse test spec: %s. Error: %v", fileName, err)
	}

	// Keep values only
	for _, command := range commandMap {
		testSpec.Commands = append(testSpec.Commands, command)
	}

	return testSpec, nil
}

func runTestSpec(t *testing.T, test *TestSpec) error {
	var cache = map[string]string{}
	for _, cmdSpec := range test.Commands {
		cmdString := generateCmdString(&cmdSpec)
		tmplOutput, tmplErr := performTemplating(strings.Join(cmdString, " "), cache)
		if tmplErr != nil {
			return fmt.Errorf("Executing templating failed: %s", tmplErr)
		}
		tmplString := strings.Fields(tmplOutput)
		t.Logf("Running: %s", strings.Join(tmplString, " "))
		actualOutput, cmdErr := exec.Command(tmplString[0], tmplString[1:]...).CombinedOutput()
		expectedOutput := regexp.MustCompile(cmdSpec.Expectation)
		if !expectedOutput.MatchString(string(actualOutput)) {
			return fmt.Errorf("miss matched expected output: %s, expectation was: %s\n", actualOutput, cmdSpec.Expectation)
		}
		if cmdErr != nil && !cmdSpec.ExpectErrorStatus {
			return fmt.Errorf("Command was expected to exit with zero status but got: %v", cmdErr)
		}
		if cmdErr == nil && cmdSpec.ExpectErrorStatus {
			return fmt.Errorf("Command was expected to exit with error status but exited with zero")
		}
	}
	return nil
}

func generateCmdString(cmdSpec *CommandSpec) (cmdString []string) {
	cmdSplit := strings.Fields(cmdSpec.Cmd)
	optionsSplit := []string{}
	for _, val := range cmdSpec.Options {
		optionsSplit = append(optionsSplit, strings.Fields(val)...)
	}
	cmdString = append(cmdSplit, cmdSpec.Args...)
	cmdString = append(cmdString, optionsSplit...)
	if regexMap[cmdSpec.Expectation] != "" {
		cmdSpec.Expectation = regexMap[cmdSpec.Expectation]
	}
	return
}

func loadRegexLookup() error {
	files, err := ioutil.ReadDir(lookupDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		err := parseLookup(path.Join(lookupDir, file.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

func parseLookup(file string) error {
	if filepath.Ext(file) != ".yml" {
		return nil
	}
	pairs, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("Unable to load regex lookup: %s. Error: %v", file, err)
	}
	if err := yaml.Unmarshal(pairs, &regexMap); err != nil {
		return fmt.Errorf("Unable to parse regex lookup: %s. Error: %v", file, err)
	}
	return nil
}

func performTemplating(s string, cache map[string]string) (output string, err error) {
	fmt.Println(s)
	var t *template.Template
	t, err = template.New("Command").Parse(s)
	if err != nil {
		return
	}
	f := func(in string) string {
		if val, ok := cache[in]; ok {
			return val
		}
		out := in + "-" + randString(10)
		cache[in] = out
		return out
	}
	var doc bytes.Buffer
	var fm = template.FuncMap{
		"uniq": func(in string) string { return f(in) },
	}
	err = t.Execute(&doc, fm)
	if err != nil {
		return
	}
	output = doc.String()
	fmt.Println(output)
	return
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
