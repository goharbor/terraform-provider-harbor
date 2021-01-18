package models

type ProjectMembersBody struct {
	ID          int                     `json:"id,omitempty"`
	RoleID      int                     `json:"role_id,omitempty"`
	GroupMember ProjectMembersBodyGroup `json:"member_group,omitempty"`
	UserMembers ProjectMemberUsersGroup `json:"member_user,omitempty"`
}

type ProjectMembersBodyGroup struct {
	GroupType   int    `json:"group_type,omitempty"`
	GroupName   string `json:"group_name,omitempty"`
	GroupID     int    `json:"id,omitempty"`
	LdapGroupDN string `json:"ldap_group_dn,omitempty"`
}

type ProjectMemberUsersGroup struct {
	UserName string `json:"username,omitempty"`
}
