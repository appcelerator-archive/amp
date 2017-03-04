package accounts

// GetOwners get teams owners
func (t *Team) GetOwners() (owners []*TeamMember) {
	for _, member := range t.Members {
		if member.Role == TeamRole_TEAM_OWNER {
			owners = append(owners, member)
		}
	}
	return
}

func (t *Team) GetMemberIndex(memberName string) int {
	memberIndex := -1
	for i, member := range t.Members {
		if member.Name == memberName {
			memberIndex = i
			break
		}
	}
	return memberIndex
}

func (o *Team) HasMember(memberName string) bool {
	return o.GetMemberIndex(memberName) != -1
}
