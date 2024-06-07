package client

import (
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func PreheatInstanceBody(d *schema.ResourceData) models.PreheatInstance {
	authInfo := models.PreheatInstanceAuthInfo{}
	authMode := d.Get("auth_mode").(string)

	switch authMode {
	case "NONE":
		// No token, username, or password
	case "BASIC":
		authInfo.Username = d.Get("username").(string)
		authInfo.Password = d.Get("password").(string)
	case "OAUTH":
		authInfo.Token = d.Get("token").(string)
	}

	body := models.PreheatInstance{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Vendor:      d.Get("vendor").(string),
		Endpoint:    d.Get("endpoint").(string),
		AuthMode:    authMode,
		AuthInfo:    authInfo,
		Enabled:     d.Get("enabled").(bool),
		Default:     d.Get("default").(bool),
		Insecure:    d.Get("insecure").(bool),
	}
	return body
}
