package schema

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

var nameFormat = regexp.MustCompile(`^[a-z0-9\-]{4,32}$`)

func isEmpty(s string) bool {
	return s == "" || strings.TrimSpace(s) == ""
}

// CheckName checks user name
func CheckName(name string) error {
	if isEmpty(name) {
		return fmt.Errorf("name is mandatory")
	}
	if !nameFormat.MatchString(name) {
		return fmt.Errorf("name is invalid")
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

func checkOrganizationMember(member *OrganizationMember) error {
	if isEmpty(member.Name) {
		return fmt.Errorf("organization member name is mandatory")
	}
	return nil
}

func checkOrganizationMembers(members []*OrganizationMember) error {
	if len(members) == 0 {
		return fmt.Errorf("organization members cannot be empty")
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
		return fmt.Errorf("organization must have at least one owner")
	}
	return nil
}

func checkTeamMember(member *TeamMember) error {
	if isEmpty(member.Name) {
		return fmt.Errorf("team member name is mandatory")
	}
	return nil
}

func checkTeamMembers(members []*TeamMember) error {
	if len(members) == 0 {
		return fmt.Errorf("team members cannot be empty")
	}
	for _, member := range members {
		if err := checkTeamMember(member); err != nil {
			return err
		}
	}
	return nil
}

func checkTeam(team *Team) (err error) {
	if err = CheckName(team.Name); err != nil {
		return err
	}
	if err = checkTeamMembers(team.Members); err != nil {
		return err
	}
	return nil
}

func checkTeams(teams []*Team) error {
	for _, team := range teams {
		if err := checkTeam(team); err != nil {
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
	if err = checkTeams(o.Teams); err != nil {
		return err
	}
	return nil
}

// Validate validates Team
func (t *Team) Validate() error {
	return checkTeam(t)
}
