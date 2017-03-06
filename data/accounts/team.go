package accounts

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
