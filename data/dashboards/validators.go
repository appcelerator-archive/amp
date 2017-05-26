package dashboards

import (
	"regexp"
	"strings"
)

var nameFormat = regexp.MustCompile(`^[a-zA-Z0-9\- ]{2,128}$`)

// CheckName checks name
func CheckName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", InvalidName
	}
	if !nameFormat.MatchString(name) {
		return "", InvalidName
	}
	return name, nil
}

// Validate validates Dashboard
func (f *Dashboard) Validate() (err error) {
	if f.Name, err = CheckName(f.Name); err != nil {
		return err
	}
	return nil
}
