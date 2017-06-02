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

func checkOrganizationMember(member *OrganizationMember) (err error) {
	if member.Name, err = CheckName(member.Name); err != nil {
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

func checkTeamMember(member string) (err error) {
	if member, err = CheckName(member); err != nil {
		return err
	}
	return nil
}

func checkTeamMembers(members []string) error {
	for _, member := range members {
		if err := checkTeamMember(member); err != nil {
			return err
		}
	}
	return nil
}

func checkTeamResource(resource *TeamResource) (err error) {
	if resource.Id, err = CheckID(resource.Id); err != nil {
		return err
	}
	return nil
}

func checkTeamResources(resources []*TeamResource) error {
	for _, resource := range resources {
		if err := checkTeamResource(resource); err != nil {
			return err
		}
	}
	return nil
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

// Validate validates Organization
func (o *Organization) Validate() (err error) {
	if o.Name, err = CheckName(o.Name); err != nil {
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
func (t *Team) Validate() (err error) {
	if t.Name, err = CheckName(t.Name); err != nil {
		return err
	}
	if err := checkTeamMembers(t.Members); err != nil {
		return err
	}
	if err := checkTeamResources(t.Resources); err != nil {
		return err
	}
	return nil
}
