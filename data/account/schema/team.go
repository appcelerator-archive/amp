package schema

// GetOwners get teams owners
func (t *Team) GetOwners() (owners []*TeamMember) {
	for _, member := range t.Members {
		if member.Role == TeamRole_TEAM_OWNER {
			owners = append(owners, member)
		}
	}
	return
}
