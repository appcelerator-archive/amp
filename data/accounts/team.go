package accounts

func (t *Team) getMemberIndex(memberName string) int {
	memberIndex := -1
	for i, member := range t.Members {
		if member == memberName {
			memberIndex = i
			break
		}
	}
	return memberIndex
}

func (t *Team) hasMember(memberName string) bool {
	return t.getMemberIndex(memberName) != -1
}

func (t *Team) getResourceIndex(resourceID string) int {
	resourceIndex := -1
	for i, resource := range t.Resources {
		if resource.Id == resourceID {
			resourceIndex = i
			break
		}
	}
	return resourceIndex
}

func (t *Team) hasResource(resourceID string) bool {
	return t.getResourceIndex(resourceID) != -1
}

func (t *Team) getResourceById(resourceID string) *TeamResource {
	for _, resource := range t.Resources {
		if resource.Id == resourceID {
			return resource
		}
	}
	return nil
}
