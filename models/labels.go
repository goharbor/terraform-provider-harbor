package models

var PathLabel = "/labels"

type Labels struct {
	Description string `json:"description"`
	Color       string `json:"color"`
	Deleted     bool   `json:"deleted"`
	Scope       string `json:"scope"`
	ProjectID   int    `json:"project_id"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
}
