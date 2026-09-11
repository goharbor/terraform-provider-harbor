package client

import (
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProjectMembersGroupBodyByID builds a project member POST/PUT body that
// attaches an already-existing Harbor usergroup by its numeric id.
//
// This is the only shape the Harbor members API handles reliably: passing
// ldap_group_dn here triggers a server-side 500 on Harbor versions that create
// the backing usergroup as a side effect. Callers must resolve the DN to an id
// first (see Client.ResolveOrCreateLdapGroup).
func ProjectMembersGroupBodyByID(role string, groupID int) models.ProjectMembersBodyPost {
	return models.ProjectMembersBodyPost{
		RoleID: RoleType(role),
		GroupMember: models.ProjectMembersBodyGroup{
			GroupID: groupID,
		},
	}
}

func ProjectMembersUserBody(d *schema.ResourceData) models.ProjectMembersBodyPost {
	return models.ProjectMembersBodyPost{
		RoleID: RoleType(d.Get("role").(string)),
		UserMembers: models.ProjectMemberUsersGroup{
			UserName: d.Get("user_name").(string),
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
		x = "maintainer"
	case 5:
		x = "limitedguest"
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
	case "maintainer":
		x = 4
	case "limitedguest":
		x = 5
	}
	return x
}
