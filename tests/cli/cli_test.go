package cli

import (
	"fmt"
	"os/exec"
	"sync"

	"errors"
	"regexp"
	"strings"
	"testing"
	"time"
)

// TestSpec contains all the CommandSpec objects
type TestSpec struct {
	Name     string
	Timeout  time.Duration
	Commands []CommandSpec
}

// CommandSpec defines the commands with arguments and options
type CommandSpec struct {
	Cmd         string   `yaml:"cmd"`
	Args        []string `yaml:"args"`
	Options     []string `yaml:"options"`
	Expectation string   `yaml:"expectation"`
	Skip        bool     `yaml:"skip"`
	Retry       int      `yaml:"retry"`
	Timeout     string   `yaml:"timeout"`
	Delay       string   `yaml:"delay"`
}

var (
	lookupDir   = "./lookup"
	setupDir    = "./setup"
	sampleDir   = "./samples"
	tearDownDir = "./tearDown"
	wg          sync.WaitGroup
	regexMap    map[string]string
)

// read, parse and execute test commands
func TestCliCmds(t *testing.T) {

	// test suite timeout
	suiteTimeout := "10m"
	duration, err := time.ParseDuration(suiteTimeout)
	if err != nil {
		t.Errorf("Unable to create duration for timeout: Suite. Error:", err)
		return
	}
	// create test suite context
	cancel := createTimeout(t, duration, "Suite")
	defer cancel()

	// parse regexes
	regexMap, err = parseLookup(lookupDir)
	if err != nil {
		t.Errorf("Unable to load lookup specs, reason:", err)
		return
	}

	// create setup timeout and parse setup specs
	setupTimeout := "8m"
	setup, err := createTestSpecs(setupDir, setupTimeout)
	if err != nil {
		t.Errorf("Unable to create setup specs, reason:", err)
		return
	}

	// create samples timeout and parse sample specs
	sampleTimeout := "30s"
	samples, err := createTestSpecs(sampleDir, sampleTimeout)
	if err != nil {
		t.Errorf("Unable to create sample specs, reason:", err)
		return
	}

	// create teardown timeout and parse tearDown specs
	tearDownTimeout := "1.5m"
	tearDown, err := createTestSpecs(tearDownDir, tearDownTimeout)
	if err != nil {
		t.Errorf("Unable to create tearDown specs, reason:", err)
		return
	}

	noOfSpecs := len(samples)

	runFramework(t, setup)
	wg.Add(noOfSpecs)
	runTests(t, samples)
	wg.Wait()
	runFramework(t, tearDown)
}

func createTestSpecs(directory string, timeout string) ([]*TestSpec, error) {
	// test spec timeout
	duration, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, fmt.Errorf("Unable to create duration for timeout: %s. Error: %v", directory, err)
	}
	// parse tests
	testSpecs, err := parseSpec(directory, duration)
	if err != nil {
		return nil, fmt.Errorf("Unable to load test specs: %s. reason: %v", directory, err)
	}
	return testSpecs, nil
}

// runs a framework (setup/tearDown)
func runFramework(t *testing.T, commands []*TestSpec) {
	for _, command := range commands {
		runFrameworkSpec(t, command)
	}
}

// runs test commands
func runTests(t *testing.T, samples []*TestSpec) {
	for _, sample := range samples {
		go runSampleSpec(t, sample)
	}
}

// execute framework commands and check for timeout, delay and retry
func runFrameworkSpec(t *testing.T, test *TestSpec) {

	var cache = map[string]string{}

	// create test spec context
	cancel := createTimeout(t, test.Timeout, test.Name)
	defer cancel()

	// iterate through all the testSpec
	for _, command := range test.Commands {
		runCmdSpec(t, command, cache)
	}
}

// execute sample commands and decrement waitgroup counter
func runSampleSpec(t *testing.T, test *TestSpec) {

	var cache = map[string]string{}

	// decrements wg counter
	defer wg.Done()

	// create test spec context
	cancel := createTimeout(t, test.Timeout, test.Name)
	defer cancel()

	// iterate through all the testSpec
	for _, command := range test.Commands {
		runCmdSpec(t, command, cache)
	}
}

//
func runCmdSpec(t *testing.T, cmd CommandSpec, cache map[string]string) {

	var i int
	var err error

	// command Spec timeout
	duration, duraErr := time.ParseDuration(cmd.Timeout)
	if duraErr != nil {
		t.Fatal("Unable to create duration for timeout:", cmd.Cmd, "Error:", duraErr)
	}
	// command Spec context
	cancel := createTimeout(t, duration, cmd.Cmd)
	defer cancel()

	// generate command slice from cmdSpec
	cmdSlice := generateCmdString(cmd)
	cmdString := strings.Join(cmdSlice, " ")

	// perform templating on command
	cmdTmplOutput, tmplErr := templating(cmdString, cache)
	if tmplErr != nil {
		t.Fatal("Executing templating failed:", cmdString, "Error:", tmplErr)
	}
	cmdTmplString := strings.Fields(cmdTmplOutput)

	// checks if the expectation has a corresponding regex
	if cmd.Expectation != "" && regexMap[cmd.Expectation] == "" {
		t.Fatal("Unable to fetch regex for command:", cmdTmplString, "reason: no regex for given expectation:", cmd.Expectation)
	}

	// perform templating on RegEx string
	regexTmplOutput, tmplErr := templating(regexMap[cmd.Expectation], cache)
	if tmplErr != nil {
		t.Fatal("Executing templating failed:", cmd.Expectation, "Error:", tmplErr)
	}

	for i = 0; i <= cmd.Retry; i++ {
		// err is set to nil a the beginning of the loop to ensure that each time a
		// command is retried or executed atleast once without the error assigned
		// from the previous executions
		err = nil

		// execute command
		cmdOutput, _ := exec.Command(cmdTmplString[0], cmdTmplString[1:]...).CombinedOutput()

		// check if the command output matches the RegEx
		expectedOutput := regexp.MustCompile(regexTmplOutput)
		if !expectedOutput.MatchString(string(cmdOutput)) {
			errString := "Mismatched expected output: " + string(cmdOutput)
			err = errors.New(errString)
		}

		// add delay to wait after command execution
		if cmd.Delay != "" {
			del, delErr := time.ParseDuration(cmd.Delay)
			if delErr != nil {
				t.Fatal("Invalid delay specified: ", cmd.Delay, "Error:", delErr)
			}
			time.Sleep(del)
		}

		// If there is no error, break the retry loop as there is no need to continue
		// If there is an error after all the retries have been used, fail the test
		if err == nil {
			break
		}
	}
	// If the command did retry, log the no. of times it did
	if i > 1 {
		t.Log("The command :", cmdTmplString, "has re-run", i, "times.")
	}
	if err != nil {
		t.Fatal(err)
	}
}

// create an array of strings representing the commands by concatenating
// all the fields from the yml files in test_samples directory
func generateCmdString(cmdSpec CommandSpec) (cmdString []string) {
	cmdSplit := strings.Fields(cmdSpec.Cmd)
	optionsSplit := []string{}
	// Possible to have multiple options
	for _, val := range cmdSpec.Options {
		optionsSplit = append(optionsSplit, strings.Fields(val)...)
	}
	cmdString = append(cmdSplit, cmdSpec.Args...)
	cmdString = append(cmdString, optionsSplit...)
	return
}
