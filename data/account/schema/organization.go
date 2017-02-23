package schema

// Validate validates Organization
func (o *Organization) GetOwners() (owners []*OrganizationMember) {
	for _, member := range o.Members {
		if member.Role == OrganizationRole_OWNER {
			owners = append(owners, member)
		}
	}
	return
}
