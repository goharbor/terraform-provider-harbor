package models

import "time"

var PathVuln = "/system/scanAll/schedule"
var PathScanners = "/scanners"

type VulnBody struct {
	Schedule struct {
		Type string `json:"type,omitempty"`
		Cron string `json:"cron,omitempty"`
	} `json:"schedule,omitempty"`
	ID            int       `json:"id,omitempty"`
	JobName       string    `json:"job_name,omitempty"`
	JobKind       string    `json:"job_kind,omitempty"`
	JobParameters string    `json:"job_parameters,omitempty"`
	JobStatus     string    `json:"job_status,omitempty"`
	Deleted       bool      `json:"deleted,omitempty"`
	CreationTime  time.Time `json:"creation_time,omitempty"`
	UpdateTime    time.Time `json:"update_time,omitempty"`
	Parameters    struct {
		AdditionalProp1 bool `json:"additionalProp1,omitempty"`
		AdditionalProp2 bool `json:"additionalProp2,omitempty"`
		AdditionalProp3 bool `json:"additionalProp3,omitempty"`
	} `json:"parameters"`
}

type ScannerBody struct {
	UUID            string    `json:"uuid,omitempty"`
	Name            string    `json:"name,omitempty"`
	Description     string    `json:"description,omitempty"`
	URL             string    `json:"url,omitempty"`
	Disabled        bool      `json:"disabled,omitempty"`
	IsDefault       bool      `json:"is_default,omitempty"`
	Auth            string    `json:"auth,omitempty"`
	SkipCertVerify  bool      `json:"skip_certVerify,omitempty"`
	UseInternalAddr bool      `json:"use_internal_addr,omitempty"`
	CreateTime      time.Time `json:"create_time,omitempty"`
	UpdateTime      time.Time `json:"update_time,omitempty"`
}
