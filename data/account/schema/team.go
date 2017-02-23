package schema

// GetOwners get teams owners
func (o *Team) GetOwners() (owners []*TeamMember) {
	for _, member := range o.Members {
		if member.Role == TeamRole_TEAM_OWNER {
			owners = append(owners, member)
		}
	}
	return
}
