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
