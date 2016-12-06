package cli

import (
	"fmt"

	"io/ioutil"
	"path"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

// read lookup directory by parsing its contents
func parseLookup(directory string) error {
	files, err := ioutil.ReadDir(lookupDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		err := generateRegexes(path.Join(lookupDir, file.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

// parse lookup directory and unmarshal its contents
func generateRegexes(fileName string) error {
	if filepath.Ext(fileName) != ".yml" {
		return fmt.Errorf("Cannot parse non-yaml file: %s", fileName)
	}
	pairs, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("Unable to read yaml regex lookup: %s. Error: %v", fileName, err)
	}
	if err := yaml.Unmarshal(pairs, &regexMap); err != nil {
		return fmt.Errorf("Unable to unmarshal yaml regex lookup: %s. Error: %v", fileName, err)
	}
	return nil
}

// read specs from directory by parsing its contents
func parseSpec(directory string, timeout time.Duration) ([]*TestSpec, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	tests := []*TestSpec{}
	for _, file := range files {
		test, err := generateTestSpecs(path.Join(directory, file.Name()), timeout)
		if err != nil {
			return nil, err
		}
		if test != nil {
			tests = append(tests, test)
		}
	}
	return tests, nil
}

// parse samples directory and unmarshal its contents
func generateTestSpecs(fileName string, timeout time.Duration) (*TestSpec, error) {
	if filepath.Ext(fileName) != ".yml" {
		return nil, fmt.Errorf("Cannot parse non-yaml file: %s", fileName)
	}
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("Unable to read yaml test spec: %s. Error: %v", fileName, err)
	}
	testSpec := &TestSpec{
		Name:    fileName,
		Timeout: timeout,
	}
	var commandMap []CommandSpec
	if err = yaml.Unmarshal(contents, &commandMap); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal yaml test spec: %s. Error: %v", fileName, err)
	}
	for _, command := range commandMap {
		if command.Timeout == "" {
			// command spec timeout
			command.Timeout = "30s"
		}
		testSpec.Commands = append(testSpec.Commands, command)
	}
	return testSpec, nil
}
