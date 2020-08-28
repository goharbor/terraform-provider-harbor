package client

import (
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func ProjectMembersBody(d *schema.ResourceData) models.ProjectMembersBody {
	return models.ProjectMembersBody{
		RoleID: RoleType(d.Get("role").(string)),
		GroupMember: models.ProjectMembersBodyGroup{
			GroupType: GroupType(d.Get("type").(string)),
			GroupName: d.Get("name").(string),
		},
	}
}

func GroupType(group string) (x int) {
	switch group {
	case "ldap":
		x = 1
	case "internal":
		x = 2
	case "oidc":
		x = 3
	}
	return x
}

func RoleTypeNumber(role int) (x string) {
	switch role {
	case 1:
		x = "projectadmin"
	case 2:
		x = "developer"
	case 3:
		x = "guest"
	case 4:
		x = "master"
	}
	return x
}

func RoleType(role string) (x int) {
	switch role {
	case "projectadmin":
		x = 1
	case "developer":
		x = 2
	case "guest":
		x = 3
	case "master":
		x = 4
	}
	return x
}
