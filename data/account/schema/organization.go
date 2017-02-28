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

// GetOwners get organization owners
func (o *Organization) IsOwner(name string) bool {
	for _, member := range o.GetMembers() {
		if member.GetName() == name {
			return true
		}
	}
	return false
}

func (o *Organization) GetMemberIndex(memberName string) int {
	memberIndex := -1
	for i, member := range o.Members {
		if member.Name == memberName {
			memberIndex = i
			break
		}
	}
	return memberIndex
}

func (o *Organization) HasMember(memberName string) bool {
	return o.GetMemberIndex(memberName) != -1
}

func (o *Organization) GetTeam(teamName string) *Team {
	for _, team := range o.Teams {
		if team.Name == teamName {
			return team
		}
	}
	return nil
}

func (o *Organization) GetTeamIndex(teamName string) int {
	teamIndex := -1
	for i, team := range o.Teams {
		if team.Name == teamName {
			teamIndex = i
			break
		}
	}
	return teamIndex
}

func (o *Organization) HasTeam(teamName string) bool {
	return o.GetTeamIndex(teamName) != -1
}
