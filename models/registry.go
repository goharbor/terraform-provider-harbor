package models

var PathRegistries = "/registries"

type RegistryBody struct {
	Status     string `json:"status,omitempty"`
	Credential struct {
		AccessKey    string `json:"access_key,omitempty"`
		AccessSecret string `json:"access_secret,omitempty"`
		Type         string `json:"type,omitempty"`
	} `json:"credential,omitempty"`
	UpdateTime   string `json:"update_time,omitempty"`
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	Insecure     bool   `json:"insecure,omitempty"`
	CreationTime string `json:"creation_time,omitempty"`
	Type         string `json:"type,omitempty"`
	ID           int    `json:"id,omitempty"`
	Description  string `json:"description,omitempty"`
}

type RegistryUpdateBody struct {
	AccessKey      string `json:"access_key"`
	CredentialType string `json:"credential_type,omitempty"`
	Name           string `json:"name,omitempty"`
	AccessSecret   string `json:"access_secret"`
	URL            string `json:"url,omitempty"`
	Insecure       bool   `json:"insecure"`
	Description    string `json:"description"`
}
