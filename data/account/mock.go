package account

import "github.com/appcelerator/amp/data/schema"

const hash = "$s2$16384$8$1$42JtddBgSqrJMwc3YuTNW+R+$ISfEF3jkvYQYk4AK/UFAxdqnmNFVeUw2gUVXEMBDAng=" // password

// Mock impliments account data.Interface
type Mock struct {
	accounts            []*schema.Account
	organizationMembers []*schema.OrganizationMember
	teams               []*schema.Team
	teamMemberships     []*schema.TeamMember
	resources           []*schema.Resource
	resourceSettings    []*schema.ResourceSettings
	permissions         []*schema.Permission
	resourceTypes       []*schema.ResourceType
	grantTypes          []*schema.GrantType
}

// NewMock returns a mock account database with some starter data
func NewMock() Interface {
	return &Mock{
		accounts: []*schema.Account{
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
		},
		organizationMembers: []*schema.OrganizationMember{
			{
				Id:            "0",
				OrgAccountId:  "1",
				UserAccountId: "0",
			},
		},
		teams: []*schema.Team{
			{
				Id:           "0",
				OrgAccountId: "1",
				Name:         "owners",
				Desc:         "Axway owners team",
			},
		},
		teamMemberships: []*schema.TeamMember{
			{
				Id:            "0",
				UserAccountId: "0",
				TeamId:        "0",
			},
		},
	}
}
