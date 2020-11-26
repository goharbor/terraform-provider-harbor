package client

import (
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetConfigSystem(d *schema.ResourceData) models.ConfigBodyPost {
	return models.ConfigBodyPost{
		ProjectCreationRestriction: d.Get("project_creation_restriction").(string),
		ReadOnly:                   d.Get("read_only").(bool),
		RobotTokenDuration:         days2mins(d.Get("robot_token_expiration").(int)),
	}
}

func GetConfigAuth(d *schema.ResourceData) models.ConfigBodyPost {
	return models.ConfigBodyPost{
		AuthMode:         d.Get("auth_mode").(string),
		OidcName:         d.Get("oidc_name").(string),
		OidcEndpoint:     d.Get("oidc_endpoint").(string),
		OidcClientID:     d.Get("oidc_client_id").(string),
		OidcClientSecret: d.Get("oidc_client_secret").(string),
		OidcGroupsClaim:  d.Get("oidc_groups_claim").(string),
		OidcScope:        d.Get("oidc_scope").(string),
		OidcVerifyCert:   d.Get("oidc_verify_cert").(bool),
		OidcAutoOnboard:  d.Get("oidc_auto_onboard").(bool),
		OidcUserClaim:    d.Get("oidc_user_claim").(string),
	}
}

func GetConfigEmail(d *schema.ResourceData) models.ConfigBodyPost {
	return models.ConfigBodyPost{
		EmailHost:     d.Get("email_host").(string),
		EmailPort:     d.Get("email_port").(int),
		EmailUsername: d.Get("email_username").(string),
		EmailPassword: d.Get("email_password").(string),
		EmailFrom:     d.Get("email_from").(string),
		EmailSsl:      d.Get("email_ssl").(bool),
		EmailInsecure: d.Get("email_insecure").(bool),
	}
}

func days2mins(days int) int {
	mins := days * 1440
	return mins
}
