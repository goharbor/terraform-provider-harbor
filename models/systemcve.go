package models
var PathSystemCVEAllowList = "/system/CVEAllowlist"
type SystemCveAllowListBodyPost struct {
	ID           int                     `json:"id,omitempty"`
	Items        SystemCveAllowlistItems `json:"items,omitempty"`
	UpdateTime   string                  `json:"update_time,omitempty"`
	CreationTime string                  `json:"creation_time,omitempty"`
	ExpiresAt    int                     `json:"expires_at,omitempty"`
}

type SystemCveAllowlistItems []struct {
	CveID string `json:"cve_id,omitempty"`
}

type SystemCveAllowlistItem struct {
	CveID string `json:"cve_id,omitempty"`
}
