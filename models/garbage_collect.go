package models

var PathGC = "/system/gc/schedule"

type GCBodyPost struct {
	Schedule struct {
		Type 			string `json:"type,omitempty"`
		Cron      string `json:"cron,omitempty"`
	} `json:"schedule,omitempty"`
	Parameters struct {
		DeleteUntagged bool `json:"delete_untagged"`
	} `json:"parameters,omitempty"`
}

type GCBodyResponses struct {
	Id int `json:"id,omitempty"`
	Schedule struct {
		Type 			string `json:"type,omitempty"`
		Cron      string `json:"cron,omitempty"`
	} `json:"schedule,omitempty"`
	JobParameters string `json:"job_parameters,omitempty"`
}