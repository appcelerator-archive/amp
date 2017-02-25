package schema

// GetOwners get organization owners
func (o *Organization) GetOwners() (owners []*OrganizationMember) {
	for _, member := range o.Members {
		if member.Role == OrganizationRole_ORGANIZATION_OWNER {
			owners = append(owners, member)
		}
	}
	return
}

// GetTeam fetches a team by name
func (o *Organization) GetTeam(name string) *Team {
	// Check if team already exists
	for _, t := range o.Teams {
		if t.Name == name {
			return t
		}
	}
	return nil
}
