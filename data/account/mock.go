package data

import "github.com/appcelerator/amp/data/schema"

const hash = "$s2$16384$8$1$42JtddBgSqrJMwc3YuTNW+R+$ISfEF3jkvYQYk4AK/UFAxdqnmNFVeUw2gUVXEMBDAng=" // password

var mockAccounts = []*schema.Account{
	{
		Id:           "0",
		Name:         "generalhenry",
		Type:         schema.AccountType_USER,
		Email:        "hallentilford@axway.com",
		PasswordHash: hash,
		IsVerified:   true,
	},
	{
		Id:         "1",
		Name:       "axway",
		Type:       schema.AccountType_ORGANIZATION,
		Email:      "hallentilford@axway.com",
		IsVerified: true,
	},
}

var mockOrganizationMembers = []*schema.OrganizationMember{
	{
		Id:            "0",
		OrgAccountId:  "1",
		UserAccountId: "0",
	},
}

var mockTeams = []*schema.Team{
	{
		Id:           "0",
		OrgAccountId: "1",
		Name:         "owners",
		Desc:         "Axway owners team",
	},
}

var mockTeamMemberships = []*schema.TeamMember{
	{
		Id:            "0",
		UserAccountId: "0",
		TeamId:        "0",
	},
}

var mockResources = []schema.Resource{}
var mockResourceSettings = []schema.ResourceSettings{}
var mockPermissions = []schema.Permission{}
var mockResourceTypes = []schema.ResourceType{}
var mockGrantTypes = []schema.GrantType{}
