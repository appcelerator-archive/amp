package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Preferences struct {
	Org  string
	Team string
}

const (
	//Preferences suffix for file
	preferences = "-preferences"
)

// SaveOrg saves the active org to file
func SaveOrg(org string, server string) error {
	existingPrefs, _ := readPreferences(server)
	if existingPrefs == nil {
		existingPrefs = &Preferences{}
	}
	existingPrefs.Org = org
	if err := savePreferences(existingPrefs, server); err != nil {
		return err
	}
	return nil
}

// SaveTeam saves the active team to file
func SaveTeam(team string, server string) error {
	existingPrefs, _ := readPreferences(server)
	if existingPrefs == nil {
		existingPrefs = &Preferences{}
	}
	existingPrefs.Team = team
	if err := savePreferences(existingPrefs, server); err != nil {
		return err
	}
	return nil
}

// savePreferences saves the preferences to file
func savePreferences(prefs *Preferences, server string) error {
	ampPrefFile := strings.TrimSuffix(server, DefaultPort) + preferences + ".yml"
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("cannot get current user: %s", err)
	}
	if err := os.MkdirAll(filepath.Join(usr.HomeDir, ampConfigFolder), os.ModePerm); err != nil {
		return fmt.Errorf("cannot create folder: %s", err)
	}
	data, err := yaml.Marshal(prefs)
	if err != nil {
		return fmt.Errorf("cannot marshal preferences: %s", err)
	}
	if err := ioutil.WriteFile(filepath.Join(usr.HomeDir, ampConfigFolder, ampPrefFile), data, 0600); err != nil {
		return fmt.Errorf("cannot write preferences: %s", err)
	}
	return nil
}

// ReadOrg reads the active org from file
func ReadOrg(server string) (string, error) {
	existingPrefs, err := readPreferences(server)
	if err != nil {
		return "", err
	}
	if existingPrefs.Org == "" {
		return "", fmt.Errorf("invalid org")
	}
	return existingPrefs.Org, nil
}

// ReadTeam reads the active team from file
func ReadTeam(server string) (string, error) {
	existingPrefs, err := readPreferences(server)
	if err != nil {
		return "", err
	}
	if existingPrefs.Team == "" {
		return "", fmt.Errorf("invalid team")
	}
	return existingPrefs.Team, nil
}

// readPreferences reads the preferences from file
func readPreferences(server string) (*Preferences, error) {
	ampPrefFile := strings.TrimSuffix(server, DefaultPort) + preferences + ".yml"
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("cannot get current user: %s", err)
	}
	data, err := ioutil.ReadFile(filepath.Join(usr.HomeDir, ampConfigFolder, ampPrefFile))
	if err != nil {
		return nil, fmt.Errorf("cannot read preferences: %s", err)
	}
	preferences := &Preferences{}
	err = yaml.Unmarshal(data, &preferences)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal preferences: %s", err)
	}
	return preferences, nil
}

// RemoveFile removes the preferences file from the .config folder
func RemoveFile(server string) error {
	ampPrefFile := strings.TrimSuffix(server, DefaultPort) + preferences + ".yml"
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("cannot get current user: %s", err)
	}
	filePath := filepath.Join(usr.HomeDir, ampConfigFolder, ampPrefFile)
	_ = os.Remove(filePath)
	return nil
}
