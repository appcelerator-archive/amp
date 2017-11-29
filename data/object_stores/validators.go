package object_stores

import (
	"regexp"
	"strings"
)

// dns like
var nameFormat = regexp.MustCompile(`^[\w.-]{3,}(?:\.[\w\.-]+)*$`)

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

// Validate validates Object store
func (f *ObjectStore) Validate() (err error) {
	if f.Name, err = CheckName(f.Name); err != nil {
		return err
	}
	return nil
}
