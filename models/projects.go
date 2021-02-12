package models

var PathProjects = "/projects"

type ProjectsBodyPost struct {
	CountLimit   int    `json:"count_limit,omitempty"`
	ProjectName  string `json:"project_name,omitempty"`
	RegistryID   int    `json:"registry_id,omitempty"`
	CveAllowlist struct {
		Items CveAllowlistItems `json:"items,omitempty"`
		// Items []struct {
		// 	CveID string `json:"cve_id,omitempty"`
		// } `json:"items,omitempty"`
		ProjectID int `json:"project_id,omitempty"`
		ID        int `json:"id,omitempty"`
		ExpiresAt int `json:"expires_at,omitempty"`
	} `json:"cve_allowlist,omitempty"`
	StorageLimit int `json:"storage_limit,omitempty"`
	Metadata     struct {
		EnableContentTrust   string `json:"enable_content_trust,omitempty"`
		AutoScan             string `json:"auto_scan,omitempty"`
		Severity             string `json:"severity,omitempty"`
		ReuseSysCveAllowlist string `json:"reuse_sys_cve_allowlist,omitempty"`
		Public               string `json:"public,omitempty"`
		PreventVul           string `json:"prevent_vul,omitempty"`
	} `json:"metadata,omitempty"`
}

type ProjectsBodyResponses struct {
	UpdateTime         string `json:"update_time"`
	OwnerName          string `json:"owner_name"`
	Name               string `json:"name"`
	Deleted            bool   `json:"deleted"`
	OwnerID            int    `json:"owner_id"`
	RepoCount          int    `json:"repo_count"`
	CreationTime       string `json:"creation_time"`
	Togglable          bool   `json:"togglable"`
	ProjectID          int    `json:"project_id"`
	CurrentUserRoleID  int    `json:"current_user_role_id"`
	CurrentUserRoleIds []int  `json:"current_user_role_ids"`
	ChartCount         int    `json:"chart_count"`
	CveAllowlist       struct {
		Items []struct {
			CveID string `json:"cve_id"`
		} `json:"items"`
		ProjectID int `json:"project_id"`
		ID        int `json:"id"`
		ExpiresAt int `json:"expires_at"`
	} `json:"cve_allowlist"`
	Metadata struct {
		EnableContentTrust   string `json:"enable_content_trust"`
		AutoScan             string `json:"auto_scan"`
		Severity             string `json:"severity"`
		ReuseSysCveAllowlist string `json:"reuse_sys_cve_allowlist"`
		Public               string `json:"public"`
		PreventVul           string `json:"prevent_vul"`
	} `json:"metadata"`
}

type CveAllowlistItems []struct {
	CveID string `json:"cve_id,omitempty"`
}
type CveAllowlistItem struct {
	CveID string `json:"cve_id,omitempty"`
}
