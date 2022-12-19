package models

type InterogationsBodyResponse struct {
	Schedule struct {
		Type string `json:"type,omitempty"`
	}
}
