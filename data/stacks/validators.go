package stacks

import (
	"regexp"
	"strings"
)

var nameFormat = regexp.MustCompile(`^[a-z0-9\-]{4,128}$`)

func isEmpty(s string) bool {
	return s == "" || strings.TrimSpace(s) == ""
}

// CheckName checks user name
func CheckName(name string) error {
	if isEmpty(name) {
		return InvalidName
	}
	if !nameFormat.MatchString(name) {
		return InvalidName
	}
	return nil
}

// Validate validates Stack
func (f *Stack) Validate() (err error) {
	if err = CheckName(f.Name); err != nil {
		return err
	}
	return nil
}
