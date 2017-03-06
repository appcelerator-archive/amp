package accounts

func (t *Team) getOwners() (owners []*TeamMember) {
	for _, member := range t.Members {
		if member.Role == TeamRole_TEAM_OWNER {
			owners = append(owners, member)
		}
	}
	return
}

func (t *Team) getMemberIndex(memberName string) int {
	memberIndex := -1
	for i, member := range t.Members {
		if member.Name == memberName {
			memberIndex = i
			break
		}
	}
	return memberIndex
}

func (t *Team) hasMember(memberName string) bool {
	return t.getMemberIndex(memberName) != -1
}
