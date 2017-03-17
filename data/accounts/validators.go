package accounts

import (
	"github.com/holys/safe"
	"net/mail"
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

// CheckEmailAddress checks email address
func CheckEmailAddress(email string) (string, error) {
	address, err := mail.ParseAddress(email)
	if err != nil {
		return "", InvalidEmail
	}
	if isEmpty(address.Address) {
		return "", InvalidEmail
	}
	return address.Address, nil
}

// CheckPassword checks password
func CheckPassword(password string) error {
	if isEmpty(password) {
		return PasswordTooWeak
	}
	safety := safe.New(8, 0, 0, safe.Simple)
	if passwordStrength := safety.Check(password); passwordStrength <= safe.Simple {
		return PasswordTooWeak
	}
	return nil
}

func checkOrganizationMember(member *OrganizationMember) error {
	if err := CheckName(member.Name); err != nil {
		return err
	}
	return nil
}

func checkOrganizationMembers(members []*OrganizationMember) error {
	if len(members) == 0 {
		return AtLeastOneOwner
	}
	haveAtLeastOneOwner := false
	for _, member := range members {
		if err := checkOrganizationMember(member); err != nil {
			return err
		}
		if member.Role == OrganizationRole_ORGANIZATION_OWNER {
			haveAtLeastOneOwner = true
		}
	}
	if !haveAtLeastOneOwner {
		return AtLeastOneOwner
	}
	return nil
}

func checkTeamMember(member *TeamMember) error {
	if err := CheckName(member.Name); err != nil {
		return err
	}
	return nil
}

func checkTeamMembers(members []*TeamMember) error {
	for _, member := range members {
		if err := checkTeamMember(member); err != nil {
			return err
		}
	}
	return nil
}

// Validate validates User
func (u *User) Validate() (err error) {
	if err = CheckName(u.Name); err != nil {
		return err
	}
	if u.Email, err = CheckEmailAddress(u.Email); err != nil {
		return err
	}
	return nil
}

// Validate validates Organization
func (o *Organization) Validate() (err error) {
	if err = CheckName(o.Name); err != nil {
		return err
	}
	if o.Email, err = CheckEmailAddress(o.Email); err != nil {
		return err
	}
	if err = checkOrganizationMembers(o.Members); err != nil {
		return err
	}
	for _, team := range o.Teams {
		if err := team.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate validates Team
func (t *Team) Validate() error {
	if err := CheckName(t.Name); err != nil {
		return err
	}
	if err := checkTeamMembers(t.Members); err != nil {
		return err
	}
	return nil
}
