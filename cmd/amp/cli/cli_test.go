package cli_test

import (
	"fmt"
	"os"
	"os/exec"

	"bytes"
	"io/ioutil"
	"math/rand"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"strconv"
	"testing"
	"text/template"
	"time"

	"github.com/appcelerator/amp/api/server"
	"gopkg.in/yaml.v2"
)

//TestSpec contains all the CommandSpec objects
type TestSpec struct {
	Name     string
	Commands []CommandSpec
}

//CommandSpec defines the commands with arguments and options
type CommandSpec struct {
	Cmd         string   `yaml:"cmd"`
	Args        []string `yaml:"args"`
	Options     []string `yaml:"options"`
	Expectation string   `yaml:"expectation"`
	Retry       int      `yaml:"retry"`
	Timeout     int64    `yaml:"timeout"`
	Delay       int64    `yaml:"delay"`
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	testDir   = "./test_samples"
	lookupDir = "./lookup"
	regexMap  map[string]string
)

//start amplifier
func TestMain(m *testing.M) {
	server.StartTestServer()
	os.Exit(m.Run())
}

//read, parse and execute test commands
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

//read test_samples directory by parsing its contents
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

//parse test_samples directory and unmarshal its contents
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

	for _, command := range commandMap {
		testSpec.Commands = append(testSpec.Commands, command)
	}
	return testSpec, nil
}

//execute commands and check for timeout, delay and retries.
func runTestSpec(t *testing.T, test *TestSpec) (err error) {
	var i int
	var cache = map[string]string{}

	for _, cmdSpec := range test.Commands {
		var tmplString []string
		startTime := time.Now().UnixNano() / 1000000

		for i = -1; i < cmdSpec.Retry; i++ {
			err = nil

			cmdString := generateCmdString(&cmdSpec)
			tmplOutput, tmplErr := performTemplating(strings.Join(cmdString, " "), cache)
			if tmplErr != nil {
				err = fmt.Errorf("Executing templating failed: %s", tmplErr)
				t.Log(err)
			}
			tmplString = strings.Fields(tmplOutput)

			t.Logf("Running: %s", strings.Join(tmplString, " "))
			cmdOutput, cmdErr := exec.Command(tmplString[0], tmplString[1:]...).CombinedOutput()
			expectedOutput := regexp.MustCompile(cmdSpec.Expectation)
			if !expectedOutput.MatchString(string(cmdOutput)) {
				err = fmt.Errorf("miss matched expected output: %s : Error: %v", cmdOutput, cmdErr)
				t.Log(err)
			}

			endTime := time.Now().UnixNano() / 1000000
			
			//timeout in Millisecond
			if cmdSpec.Timeout != 0 && endTime-startTime >= cmdSpec.Timeout {
				return fmt.Errorf("Command execution has exceeded timeout : %s", tmplString)
			}
			if err == nil {
				break
			}
			//delay in Millisecond
			time.Sleep(time.Duration(cmdSpec.Delay) * time.Millisecond)
		}
		if i > 0 && i == cmdSpec.Retry {
			t.Log("This command :", tmplString, "has re-run", i, "times.")
		}
	}
	return err
}

//create an array of strings representing the commands by concatenating
//all the fields from the yml files in test_samples directory
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
	p := func(in string, min, max int) string {
		if val, ok := cache[in]; ok {
			return val
		}
		out := strconv.Itoa(rand.Intn(max - min) + min)
		cache[in] = out
		return out
	}
	var doc bytes.Buffer
	var fm = template.FuncMap{
		"uniq": func(in string) string { return f(in) },
		"port": func(in string, min, max int) string { return p(in, min, max) },
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
