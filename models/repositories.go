package models

type RepositoryBody struct {
	ArtifactCount int    `json:"artifact_count,omitempty"`
	ID            int    `json:"id,omitempty"`
	ProjectID     int    `json:"project_id,omitempty"`
	PullCount     int    `json:"pull_count,omitempty"`
	Name          string `json:"name,omitempty"`
	CreationTime  string `json:"creation_time,omitempty"`
	UpdateTime    string `json:"update_time,omitempty"`
}
