package account

import "github.com/appcelerator/amp/data/account/schema"

func FromSchema(schema *schema.User) *User {
	return &User{
		Name:     schema.Name,
		Email:    schema.Email,
		CreateDt: schema.CreateDt,
	}
}
