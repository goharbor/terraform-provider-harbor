package models

var PathGroups = "/usergroups"

//
type GroupBody struct {
	Groupname   string `json:"group_name,omitempty"`
	GroupType   int    `json:"group_type,omitempty"`
	LdapGroupDn string `json:"ldap_group_dn,omitempty"`
	ID          int    `json:"id,omitempty"`
}
