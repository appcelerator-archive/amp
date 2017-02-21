package cli

import (
	"errors"
	"fmt"
	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/docker/distribution/context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	gexpect "github.com/Thomasrooney/gexpect"
)

// SuiteSpec defines the suite fields and the suite data structures.
type SuiteSpec struct {
	Name        string   `yaml:"name"`
	Timeout     string   `yaml:"timeout"`
	LookupDir   string   `yaml:"lookupdir"`
	SetupDir    string   `yaml:"setupdir"`
	TestDirs    []string `yaml:"testdirs"`
	TearDownDir string   `yaml:"teardowndir"`
	Setup       []TestSpec
	Tests       []TestSpec
	TearDown    []TestSpec
}

// TestSpec defines the test fields and test data structures.
type TestSpec struct {
	Name        string        `yaml:"name"`
	TestTimeout string        `yaml:"test-timeout"`
	CmdTimeout  string        `yaml:"cmd-timeout"`
	Commands    []CommandSpec `yaml:"commands"`
	Cache       map[string]string
}

// CommandSpec defines the command fields and data types.
type CommandSpec struct {
	Cmd         string   `yaml:"cmd"`
	Args        []string `yaml:"args"`
	Options     []string `yaml:"options"`
	Input       []string `yaml:"input"`
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
	ctx := context.Background()

	// Connect to amplifier
	conn, err := grpc.Dial("localhost:8080",
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(60*time.Second),
	)
	assert.NoError(t, err)

	accountClient := account.NewAccountClient(conn)

	signUpRequest := account.SignUpRequest{
		Name:     "cli",
		Password: "cliPassword",
		Email:    "cli@amp.io",
	}

	// SignUp
	accountClient.SignUp(ctx, &signUpRequest)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify
	accountClient.Verify(ctx, &account.VerificationRequest{Token: token})

	// Login, somehow
	md := metadata.Pairs(auth.TokenKey, token)
	cli.SaveToken(md)

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

// runCmdSpec creates the CommandSpec timeout and executes the command.
func runCmdSpec(t *testing.T, cmdSpec CommandSpec, cache map[string]string) error {
	// Retry counter.
	i := 0

	// Cmd output.
	output := ""

	// Generates the command string.
	cmd := genCmdStr(cmdSpec)

	// Create commandspec duration.
	duration, err := time.ParseDuration(cmdSpec.Timeout)
	if err != nil {
		return fmt.Errorf("Unable to create duration for timeout:", cmd, "Error:", err)
	}
	// Create commandspec timeout.
	cancel := createTimeout(t, duration, cmd)
	defer cancel()

	// Template command string.
	cmd, err = templating(cmd, cache)
	if err != nil {
		return fmt.Errorf("Executing templating failed:", cmd, "Error:", err)
	}

	// Check if expectation has corresponding regex
	if cmdSpec.Expectation != "" && regexMap[cmdSpec.Expectation] == "" {
		return fmt.Errorf("Unable to fetch regex for command:", cmd, "reason: no regex for given expectation:", cmdSpec.Expectation)
	}

	// Template regex.
	regex, err := templating(regexMap[cmdSpec.Expectation], cache)
	if err != nil {
		return fmt.Errorf("Executing templating failed:", cmdSpec.Expectation, "Error:", err)
	}

	// Retry loop.
	for i = 0; i <= cmdSpec.Retry; i++ {
		// Execute command.
		output = runCmd(cmd, cmdSpec)

		// Check ouptput against corresponding regex.
		expectedOutput := regexp.MustCompile(regex)
		if expectedOutput.MatchString(string(output)) {
			// Passed
			return nil
		}

		// Delay command execution if delay is set.
		if cmdSpec.Delay != "" {
			del, err := time.ParseDuration(cmdSpec.Delay)
			if err != nil {
				return fmt.Errorf("Invalid delay specified: ", cmdSpec.Delay, "Error:", err)
			}
			time.Sleep(del)
		}
	}
	// Command fails as there is a mismatched output.
	return fmt.Errorf("Error: mismatched return: command, " + cmd + ". regex, " + regex + ". ouput, " + output)
}

// runCmd executes the command, passing any inputs to it.
func runCmd(cmd string, spec CommandSpec) string {
	// Spawns new execution of command.
	child, err := gexpect.Spawn(cmd)
	if err != nil {
		panic(err)
	}

	// Sends each input to command.
	for _, input := range spec.Input {
		child.Send(input + "\r")
		child.Expect(input)
	}

	// Reads ouput until return on cli.
	output, _ := child.ReadUntil('$')

	return string(output)
}

// genCmdStr creates the command by concatenating the cmd, args and options fields.
func genCmdStr(cmdSpec CommandSpec) string {
	// Get command string.
	cmdStr := cmdSpec.Cmd

	// Add arguments.
	for _, arg := range cmdSpec.Args {
		cmdStr = cmdStr + " " + arg
	}

	// Add options.
	for _, opts := range cmdSpec.Options {
		cmdStr = cmdStr + " " + opts
	}

	return cmdStr
}
