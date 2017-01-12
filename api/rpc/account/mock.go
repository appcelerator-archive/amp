package account

import "github.com/appcelerator/amp/data/schema"

var accounts = []*schema.Account{
	{
		Id:         "0",
		Name:       "generalhenry",
		Type:       schema.AccountType_USER,
		Email:      "hallentilford@axway.com",
		PwHashcode: hash,
		IsVerified: true,
	},
	{
		Id:         "1",
		Name:       "axway",
		Type:       schema.AccountType_ORGANIZATION,
		Email:      "hallentilford@axway.com",
		IsVerified: true,
	},
}

var organizationMembers = []*schema.OrganizationMember{
	{
		Id:            "0",
		OrgAccountId:  "1",
		UserAccountId: "0",
	},
}

var teams = []*schema.Team{
	{
		Id:           "0",
		OrgAccountId: "1",
		Name:         "owners",
		Desc:         "Axway owners team",
	},
}

var teamMemberships = []*schema.TeamMember{
	{
		Id:            "0",
		UserAccountId: "0",
		TeamId:        "0",
	},
}

var resources = []schema.Resource{}
var resourceSettings = []schema.ResourceSettings{}
var permissions = []schema.Permission{}
var resourceTypes = []schema.ResourceType{}
var grantTypes = []schema.GrantType{}
