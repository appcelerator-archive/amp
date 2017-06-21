package accounts

import (
	"net/mail"
	"regexp"
	"strings"

	"github.com/holys/safe"
)

var nameFormat = regexp.MustCompile(`^[a-z0-9\-]{2,128}$`)

// CheckName checks user name
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

// CheckEmailAddress checks email address
func CheckEmailAddress(email string) (string, error) {
	address, err := mail.ParseAddress(strings.TrimSpace(email))
	if err != nil {
		return "", InvalidEmail
	}
	if address.Address == "" {
		return "", InvalidEmail
	}
	return address.Address, nil
}

// CheckPassword checks password
func CheckPassword(password string) (string, error) {
	safety := safe.New(8, 0, 0, safe.Simple)
	if passwordStrength := safety.Check(password); passwordStrength <= safe.Simple {
		return "", PasswordTooWeak
	}
	return password, nil
}

// CheckID checks resource id
func CheckID(ID string) (string, error) {
	ID = strings.TrimSpace(ID)
	if ID == "" {
		return "", InvalidResourceID
	}
	return ID, nil
}

// Validate validates User
func (u *User) Validate() (err error) {
	if u.Name, err = CheckName(u.Name); err != nil {
		return err
	}
	if u.Email, err = CheckEmailAddress(u.Email); err != nil {
		return err
	}
	return nil
}
