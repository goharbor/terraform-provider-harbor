package models

type ProjectMembersBody struct {
	ID          int                     `json:"id,omitempty"`
	RoleID      int                     `json:"role_id,omitempty"`
	GroupMember ProjectMembersBodyGroup `json:"member_group,omitempty"`
}

type ProjectMembersBodyGroup struct {
	GroupType int    `json:"group_type,omitempty"`
	GroupName string `json:"group_name,omitempty"`
}
