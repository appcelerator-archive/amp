package schema

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

var nameFormat = regexp.MustCompile(`^[a-z0-9]+$`)

func isEmpty(s string) bool {
	return s == "" || strings.TrimSpace(s) == ""
}

// CheckUserName checks user name
func CheckName(name string) error {
	if isEmpty(name) {
		return fmt.Errorf("name is mandatory")
	}
	if !nameFormat.MatchString(name) {
		return fmt.Errorf("name is invlaid")
	}
	return nil
}

func checkPasswordHash(passwordHash string) error {
	if isEmpty(passwordHash) {
		return fmt.Errorf("password hash is mandatory")
	}
	return nil
}

// CheckEmailAddress checks email address
func CheckEmailAddress(email string) (string, error) {
	address, err := mail.ParseAddress(email)
	if err != nil {
		return "", err
	}
	if isEmpty(address.Address) {
		return "", fmt.Errorf("email is mandatory")
	}
	return address.Address, nil
}

// Validate validates User
func (u *User) Validate() (err error) {
	if err = CheckName(u.Name); err != nil {
		return err
	}
	if u.Email, err = CheckEmailAddress(u.Email); err != nil {
		return err
	}
	if err = checkPasswordHash(u.PasswordHash); err != nil {
		return err
	}
	return nil
}
