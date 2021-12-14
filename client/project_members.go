package client

import (
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ProjectMembersGroupBody(d *schema.ResourceData) models.ProjectMembersBodyPost {
	body := models.ProjectMembersBodyPost{
		RoleID: RoleType(d.Get("role").(string)),
		GroupMember: models.ProjectMembersBodyGroup{
			GroupType: GroupType(d.Get("type").(string)),
			GroupName: d.Get("group_name").(string),
			GroupID:   d.Get("group_id").(int),
		},
	}

	if v, ok := d.GetOk("ldap_group_dn"); ok {
		body.GroupMember.LdapGroupDN = v.(string)
	}
	return body
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
		x = "master"
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
	case "master", "maintainer":
		x = 4
	case "limitedguest":
		x = 5
	}
	return x
}
