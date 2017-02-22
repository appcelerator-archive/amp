package cli

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// createSuite reads every suite directory and gets their contents.
func createSuite(dirs []string) ([]*SuiteSpec, error) {
	// Array of suites.
	suites := []*SuiteSpec{}
	for _, dir := range dirs {

		// Read suite yaml file contents.
		contents, err := ioutil.ReadFile(path.Join(dir, "/suite.yml"))
		if err != nil {
			return nil, err
		}

		// Unmarshal contents into suitespec struct.
		suite := &SuiteSpec{}
		if err = yaml.Unmarshal(contents, suite); err != nil {
			return nil, fmt.Errorf("Unable to unmarshal yaml suite spec: %s. Error: %v", dir, err)
		}

		// Parse lookup files from given lookup directory in suitespec.
		regexMap, err = parseLookup(path.Join(dir, suite.LookupDir))
		if err != nil {
			return nil, fmt.Errorf("Unable to parse yaml lookup spec: %s. Error: %v", suite.LookupDir, err)
		}

		// Parse setup files from given setup directory in suitespec.
		setup, err := parseSpec(path.Join(dir, suite.SetupDir))
		if err != nil {
			return nil, fmt.Errorf("Unable to parse yaml setup spec: %s. Error: %v", suite.SetupDir, err)
		}
		suite.Setup = append(suite.Setup, *setup...)

		// Parse test files from given test directories in suitespec.
		for _, testDir := range suite.TestDirs {
			test, err := parseSpec(path.Join(dir, testDir))
			if err != nil {
				return nil, fmt.Errorf("Unable to parse yaml test spec: %s. Error: %v", testDir, err)
			}
			suite.Tests = append(suite.Tests, *test...)
		}

		// Parse teardown files from given teardown test directory in suitespec.
		tearDown, err := parseSpec(path.Join(dir, suite.TearDownDir))
		if err != nil {
			return nil, fmt.Errorf("Unable to parse yaml teardown spec: %s. Error: %v", suite.TearDownDir, err)
		}
		suite.TearDown = append(suite.TearDown, *tearDown...)

		suites = append(suites, suite)
	}

	return suites, nil
}

// parseLookup reads the lookup directory and gets it contents.
func parseLookup(dir string) (map[string]string, error) {
	// Read lookup directory for its files.
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Insert lookup data into Regex Map.
	rgxMap := make(map[string]string)
	for _, file := range files {
		lookup, err := createLookup(path.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		for expect, rgx := range lookup {
			rgxMap[expect] = rgx
		}
	}

	return rgxMap, nil
}

// createLookup reads the lookup files contents and unmarshals it.
func createLookup(name string) (map[string]string, error) {
	// File must be a yaml file.
	if filepath.Ext(name) != ".yml" {
		return nil, fmt.Errorf("Cannot parse non-yaml file: %s", name)
	}

	// Read lookup yaml file contents.
	contents, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("Unable to read yaml regex lookup: %s. Error: %v", name, err)
	}

	// Unmarshals lookup yaml file contents into map.
	lookup := make(map[string]string)
	if err := yaml.Unmarshal(contents, &lookup); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal yaml lookup: %s. Error: %v", name, err)
	}

	return lookup, nil
}

// parseSpec reads the given testspec directory and gets it contents.
func parseSpec(dir string) (*[]TestSpec, error) {
	// Read testspec directory for its files.
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Insert testspec data into struct.
	tests := []TestSpec{}
	for _, file := range files {
		test, err := createSpec(path.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		if test != nil {
			tests = append(tests, *test)
		}
	}

	return &tests, nil
}

// createLookup reads the testspec files contents and unmarshals it.
func createSpec(name string) (*TestSpec, error) {
	// File must be a yaml file.
	if filepath.Ext(name) != ".yml" {
		return nil, fmt.Errorf("Cannot parse non-yaml file: %s", name)
	}

	// Read testspec yaml file contents.
	contents, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("Unable to read yaml test spec: %s. Error: %v", name, err)
	}

	// Unmarshals testspec yaml file contents into struct.
	test := &TestSpec{}
	if err = yaml.Unmarshal(contents, &test); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal yaml test spec: %s. Error: %v", name, err)
	}

	// Instanstiates cache for templating.
	test.Cache = make(map[string]string)

	// Assigns default values for commands.
	for i := range test.Commands {
		// Skip command by removing from command list.
		if test.Commands[i].Skip == true {
			test.Commands = append(test.Commands[:i], test.Commands[i+1:]...)
		}

		// Default commandspec timeout.
		if test.Commands[i].Timeout == "" {
			test.Commands[i].Timeout = test.CmdTimeout
		}
	}
	return test, nil
}
