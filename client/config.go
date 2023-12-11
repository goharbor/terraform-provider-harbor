package client

import (
	"encoding/json"
	"log"

	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetConfigSystem(d *schema.ResourceData) models.ConfigBodySystemPost {
	var body models.ConfigBodySystemPost
	storage := d.Get("storage_per_project").(int)
	if storage > 0 {
		storage *= 1073741824 // GB to Byte
	}
	body = models.ConfigBodySystemPost{
		ProjectCreationRestriction: d.Get("project_creation_restriction").(string),
		ReadOnly:                   d.Get("read_only").(bool),
		ScannerSkipUpdatePulltime:  d.Get("scanner_skip_update_pulltime").(bool),
		RobotTokenDuration:         d.Get("robot_token_expiration").(int),
		QuotaPerProjectEnable:      true,
		RobotNamePrefix:            d.Get("robot_name_prefix").(string),
		StoragePerProject:          storage,
	}
	log.Printf("[DEBUG] %+v\n ", body)
	return body
}

func GetConfigAuth(d *schema.ResourceData) models.ConfigBodyAuthPost {
	var body models.ConfigBodyAuthPost

	auth := d.Get("auth_mode").(string)

	switch auth {
	case "oidc_auth", "oidc":
		body = models.ConfigBodyAuthPost{
			AuthMode:         "oidc_auth",
			PrimaryAuthMode:  d.Get("primary_auth_mode").(bool),
			OidcName:         d.Get("oidc_name").(string),
			OidcEndpoint:     d.Get("oidc_endpoint").(string),
			OidcClientID:     d.Get("oidc_client_id").(string),
			OidcClientSecret: d.Get("oidc_client_secret").(string),
			OidcGroupFilter:  d.Get("oidc_group_filter").(string),
			OidcGroupsClaim:  d.Get("oidc_groups_claim").(string),
			OidcScope:        d.Get("oidc_scope").(string),
			OidcVerifyCert:   d.Get("oidc_verify_cert").(bool),
			OidcAutoOnboard:  d.Get("oidc_auto_onboard").(bool),
			OidcUserClaim:    d.Get("oidc_user_claim").(string),
			OidcAdminGroup:   d.Get("oidc_admin_group").(string),
		}
	case "ldap_auth", "ldap":
		body = models.ConfigBodyAuthPost{
			AuthMode:                     "ldap_auth",
			PrimaryAuthMode:              d.Get("primary_auth_mode").(bool),
			LdapURL:                      d.Get("ldap_url").(string),
			LdapSearchDn:                 d.Get("ldap_search_dn").(string),
			LdapSearchPassword:           d.Get("ldap_search_password").(string),
			LdapBaseDn:                   d.Get("ldap_base_dn").(string),
			LdapFilter:                   d.Get("ldap_filter").(string),
			LdapUID:                      d.Get("ldap_uid").(string),
			LdapGroupBaseDn:              d.Get("ldap_group_base_dn").(string),
			LdapGroupSearchFilter:        d.Get("ldap_group_filter").(string),
			LdapGroupAttributeName:       d.Get("ldap_group_gid").(string),
			LdapGroupAdminDn:             d.Get("ldap_group_admin_dn").(string),
			LdapGroupMembershipAttribute: d.Get("ldap_group_membership").(string),
			LdapVerifyCert:               d.Get("ldap_verify_cert").(bool),
			LdapScope:                    getLdapScope(d.Get("ldap_scope").(string)),
			LdapGroupSearchScope:         getLdapScope(d.Get("ldap_group_scope").(string)),
		}
	}
	log.Printf("[DEBUG] %+v\n ", body)
	return body
}

func GetConfigEmail(d *schema.ResourceData) models.ConfigBodyEmailPost {
	var body models.ConfigBodyEmailPost
	body = models.ConfigBodyEmailPost{
		EmailHost:     d.Get("email_host").(string),
		EmailPort:     d.Get("email_port").(int),
		EmailUsername: d.Get("email_username").(string),
		EmailPassword: d.Get("email_password").(string),
		EmailFrom:     d.Get("email_from").(string),
		EmailSsl:      d.Get("email_ssl").(bool),
		EmailInsecure: d.Get("email_insecure").(bool),
	}
	log.Printf("[DEBUG] %+v\n ", body)
	return body
}

func days2mins(days int) int {
	mins := days * 1440
	return mins
}

func getLdapScope(scopeName string) (scopeInt int) {
	var scope int
	switch scopeName {
	case "base":
		scope = 0
	case "onelevel":
		scope = 1
	case "subtree":
		scope = 2
	}
	return scope
}
func getLdapScopeName(scopeInt int) (scopeName string) {
	switch scopeInt {
	case 0:
		scopeName = "base"
	case 1:
		scopeName = "onelevel"
	case 2:
		scopeName = "subtree"
	default:
		scopeName = "subtree"
	}
	return scopeName
}

// SetAuthValues the values for Auth when performing a read on resource
func SetAuthValues(d *schema.ResourceData, resp string) error {
	var jsonData models.ConfigBodyResponse
	err := json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return err
	}

	auth := jsonData.AuthMode.Value
	d.Set("auth_mode", auth)
	d.Set("primary_auth_mode", jsonData.PrimaryAuthMode.Value)

	switch auth {
	case "oidc_auth", "oidc":
		d.Set("oidc_name", jsonData.OidcName.Value)
		d.Set("oidc_endpoint", jsonData.OidcEndpoint.Value)
		d.Set("oidc_client_id", jsonData.OidcClientID.Value)
		d.Set("oidc_groups_claim", jsonData.OidcGroupsClaim.Value)
		d.Set("oidc_scope", jsonData.OidcScope.Value)
		d.Set("oidc_verify_cert", jsonData.OidcVerifyCert.Value)
		d.Set("oidc_auto_onboard", jsonData.OidcAutoOnboard.Value)
		d.Set("oidc_user_claim", jsonData.OidcUserClaim.Value)
		d.Set("oidc_admin_group", jsonData.OidcAdminGroup.Value)
	case "ldap_auth", "ldap":
		d.Set("ldap_url", jsonData.LdapURL.Value)
		d.Set("ldap_base_dn", jsonData.LdapBaseDn.Value)
		d.Set("ldap_uid", jsonData.LdapUID.Value)
		d.Set("ldap_verify_cert", jsonData.LdapVerifyCert.Value)
		d.Set("ldap_search_dn", jsonData.LdapSearchDn.Value)
		d.Set("ldap_group_base_dn", jsonData.LdapGroupBaseDn.Value)
		d.Set("ldap_group_filter", jsonData.LdapGroupSearchFilter.Value)
		d.Set("ldap_group_gid", jsonData.LdapGroupAttributeName.Value)
		d.Set("ldap_group_admin_dn", jsonData.LdapGroupAdminDn.Value)
		d.Set("ldap_group_membership", jsonData.LdapGroupMembershipAttribute.Value)
		d.Set("ldap_scope", getLdapScopeName(jsonData.LdapScope.Value))
		d.Set("ldap_group_scope", getLdapScopeName(jsonData.LdapGroupSearchScope.Value))
	}

	return nil
}
