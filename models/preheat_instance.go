package models

var PathPreheatInstance = "/p2p/preheat/instances"

type PreheatInstance struct {
	ID             int                     `json:"id"`
	Name           string                  `json:"name"`
	Description    string                  `json:"description"`
	Vendor         string                  `json:"vendor"`
	Endpoint       string                  `json:"endpoint"`
	AuthMode       string                  `json:"auth_mode"`
	AuthInfo       PreheatInstanceAuthInfo `json:"auth_info"`
	Status         string                  `json:"status"`
	Enabled        bool                    `json:"enabled"`
	Default        bool                    `json:"default"`
	Insecure       bool                    `json:"insecure"`
	SetupTimestamp int64                   `json:"setup_timestamp"`
}

type PreheatInstanceAuthInfo struct {
	Token    string `json:"token,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
