package models

var PathUsers = "/users"

//
type UserBody struct {
	Location      string `json:"location,omitempty"`
	Username        string `json:"username,omitempty"`
	Comment         string `json:"comment,omitempty"`
	UpdateTime      string `json:"update_time,omitempty"`
	Password        string `json:"password,omitempty"`
	UserID          int    `json:"user_id,omitempty"`
	Realname        string `json:"realname,omitempty"`
	Deleted         bool   `json:"deleted,omitempty"`
	CreationTime    string `json:"creation_time,omitempty"`
	AdminRoleInAuth bool   `json:"admin_role_in_auth,omitempty"`
	RoleID          int    `json:"role_id,omitempty"`
	SysadminFlag    bool   `json:"sysadmin_flag,omitempty"`
	RoleName        string `json:"role_name,omitempty"`
	ResetUUID       string `json:"reset_uuid,omitempty"`
	Salt            string `json:"Salt,omitempty"`
	Email           string `json:"email,omitempty"`
	Newpassword     string `json:"new_password,omitempty"`
}
