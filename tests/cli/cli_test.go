package cli

import (
	"fmt"
	"os/exec"
	"sync"

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
	Retry       int      `yaml:"retry"`
	Timeout     string   `yaml:"timeout"`
	Delay       string   `yaml:"delay"`
}

var (
	suiteTimeout string
	lookupDir    = "./lookup"
	sampleDir    = "./samples"
	regexMap     map[string]string
	wg           sync.WaitGroup
)

// read, parse and execute test commands
func TestCmds(t *testing.T) {
	// test suite timeout
	suiteTimeout = "120s"
	duration, err := time.ParseDuration(suiteTimeout)
	if err != nil {
		t.Errorf("Unable to create duration for timeout: Suite. Error: %v", err)
		return
	}
	// create test suite context
	cancelSuite := createTimeout(t, duration, "Suite")
	defer cancelSuite()
	// parse regexes
	err = parseLookup(lookupDir)
	if err != nil {
		t.Errorf("Unable to load lookup specs, reason: %v", err)
		return
	}
	// parse test samples
	tests, err := parseSpec(sampleDir, duration)
	if err != nil {
		t.Errorf("Unable to load test specs, reason: %v", err)
		return
	}
	wg.Add(len(tests))
	t.Logf("-----------------------------------------------------------------------------------------")
	for _, test := range tests {
		t.Logf("Running spec: %s", test.Name)
		go runTestSpec(t, test)
	}
	wg.Wait()
	t.Logf("-----------------------------------------------------------------------------------------")
	t.Log("Finished!")
}

// execute commands and check for timeout, delay and retry
func runTestSpec(t *testing.T, test *TestSpec) {
	defer wg.Done()
	// test spec timeout
	testSpecTimeout := "60s"
	duration, duraErr := time.ParseDuration(testSpecTimeout)
	if duraErr != nil {
		t.Errorf("Unable to create duration for timeout: TestSpec. Error: %v", duraErr)
		return
	}
	// create test spec context
	cancelTestSpec := createTimeout(t, duration, test.Name)
	defer cancelTestSpec()
	var i int
	var cache = map[string]string{}
	var err error
	// iterate through all the testSpec
	for _, cmdSpec := range test.Commands {
		var cmdTmplString []string
		duration, duraErr = time.ParseDuration(cmdSpec.Timeout)
		if duraErr != nil {
			t.Log("Unable to create duration for timeout: %s. Error: %v", cmdSpec.Cmd, duraErr)
			t.Fatal()
		}
		// cmd Spec context
		cancelCmdSpec := createTimeout(t, duration, cmdSpec.Cmd)
		for i = -1; i < cmdSpec.Retry; i++ {
			// err is set to nil a the beginning of the loop to ensure that each time a
			// command is retried or executed atleast once without the error assigned
			// from the previous executions
			err = nil

			// generate command string from cmdSpec
			cmdString := generateCmdString(&cmdSpec)

			// perform templating on cmdString
			cmdTmplOutput, tmplErr := templating(strings.Join(cmdString, " "), cache)
			if tmplErr != nil {
				t.Log("Executing templating failed: %s. Error: %v", cmdString, tmplErr)
				t.Fatal()
			}
			cmdTmplString = strings.Fields(cmdTmplOutput)

			// execute command
			cmdOutput, cmdErr := exec.Command(cmdTmplString[0], cmdTmplString[1:]...).CombinedOutput()
			t.Logf("Running: %s", strings.Join(cmdTmplString, " "))

			//perform templating on RegEx string
			regexTmplOutput, tmplErr := templating(regexMap[cmdSpec.Expectation], cache)
			if tmplErr != nil {
				t.Log("Executing templating failed: %s. Error: %v", cmdSpec.Expectation, tmplErr)
				t.Fatal()
			}

			// check if the command output matches the RegEx
			expectedOutput := regexp.MustCompile(regexTmplOutput)
			if !expectedOutput.MatchString(string(cmdOutput)) {
				err = fmt.Errorf("Mismatched expected output: %s : Error: %v", cmdOutput, cmdErr)
				t.Log(err)
			}

			// if no error after retries, break the loop to continue command execution
			if err == nil {
				break
			}
			// add delay (in Millisecond) to wait for command execution
			if cmdSpec.Delay != "" {
				del, delErr := time.ParseDuration(cmdSpec.Delay)
				if delErr != nil {
					t.Log("Invalid delay specified: %s : Error: %v", cmdSpec.Delay, delErr)
					t.Fatal()
				}
				time.Sleep(del)
			}
		}
		if i > 0 {
			t.Log("This command :", cmdTmplString, "has re-run", i, "times.")
		}
		if err != nil {
			t.Log(err)
			t.Fatal()
		}
		cancelCmdSpec()
	}
}

// create an array of strings representing the commands by concatenating
// all the fields from the yml files in test_samples directory
func generateCmdString(cmdSpec *CommandSpec) (cmdString []string) {
	cmdSplit := strings.Fields(cmdSpec.Cmd)
	optionsSplit := []string{}
	for _, val := range cmdSpec.Options {
		optionsSplit = append(optionsSplit, strings.Fields(val)...)
	}
	cmdString = append(cmdSplit, cmdSpec.Args...)
	cmdString = append(cmdString, optionsSplit...)
	return
}
