package models

type RobotBody struct {
	Access      []RobotBodyAccess `json:"access,omitempty"`
	Name        string            `json:"name,omitempty"`
	ExpiresAt   int               `json:"expires_at,omitempty"`
	Description string            `json:"description,omitempty"`
}
type RobotBodyAccess struct {
	Action   string `json:"action,omitempty"`
	Resource string `json:"resource,omitempty"`
}
type RobotBodyRepones struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Token       string `json:"token"`
	Secret      string `json:"secret"`
	Description string `json:"description"`
	ProjectID   int    `json:"project_id"`
	ExpiresAt   int    `json:"expires_at"`
	Disabled    bool   `json:"disabled"`
}
