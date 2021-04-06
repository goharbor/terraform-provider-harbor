package models

type RobotBodyPermission struct {
	Access    []RobotBodyAccess `json:"access,omitempty"`
	Kind      string            `json:"kind,omitempty"`
	Namespace string            `json:"namespace,omitempty"`
}
type RobotBodyAccess struct {
	Action   string `json:"action,omitempty"`
	Resource string `json:"resource,omitempty"`
	Effect   string `json:"effect,omitempty"`
}
type RobotBody struct {
	ID          int                   `json:"id,omitempty"`
	Name        string                `json:"name,omitempty"`
	Level       string                `json:"level,omitempty"`
	Description string                `json:"description,omitempty"`
	Secret      string                `json:"secret,omitempty"`
	Duration    int                   `json:"duration,omitempty"`
	Disable     bool                  `json:"disable,omitempty"`
	Permissions []RobotBodyPermission `json:"permissions,omitempty"`
}
type RobotBodyResponse struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Secret       string `json:"secret,omitempty"`
	ExpiresAt    int    `json:"expires_at,omitempty"`
	CreationTime string `json:"creation_time,omitempty"`
}
