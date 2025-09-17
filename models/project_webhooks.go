package models

import "time"

type ProjectWebhook struct {
	UpdateTime   time.Time        `json:"update_time"`
	Description  string           `json:"description"`
	Creator      string           `json:"creator"`
	CreationTime time.Time        `json:"creation_time"`
	Enabled      bool             `json:"enabled"`
	Targets      []WebHookTargets `json:"targets"`
	EventTypes   []interface{}    `json:"event_types"`
	ProjectID    int              `json:"project_id"`
	ID           int              `json:"id"`
	Name         string           `json:"name"`
}
type WebHookTargets struct {
	Type           string `json:"type"`
	AuthHeader     string `json:"auth_header"`
	SkipCertVerify bool   `json:"skip_cert_verify"`
	Address        string `json:"address"`
	PayloadFormat  string `json:"payload_format"`
}
