package functions

import (
	"regexp"
	"strings"
)

var nameFormat = regexp.MustCompile(`^[a-z0-9\-]{3,128}$`)

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

// CheckImage checks image name
func CheckImage(image string) (string, error) {
	image = strings.TrimSpace(image)
	if image == "" {
		return "", InvalidImage
	}
	return image, nil
}

// Validate validates Function
func (f *Function) Validate() (err error) {
	if f.Name, err = CheckName(f.Name); err != nil {
		return err
	}
	if f.Image, err = CheckImage(f.Image); err != nil {
		return err
	}
	return nil
}
