package provider

import (
	"encoding/json"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConfigAuth() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"auth_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"oidc_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oidc_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oidc_client_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oidc_client_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"oidc_groups_claim": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oidc_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oidc_verify_cert": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"oidc_auto_onboard": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"oidc_user_claim": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ldap_url": {
				Type:     schema.TypeString,
				Optional: true,
				// ConflictsWith: []string{"oidc_name", "oidc_endpoint", "oidc_client_id", "oidc_client_secret", "oidc_groups_claim", "oidc_scope", "oidc_verify_cert", "oidc_auto_onboard", "oidc_user_claim"},
				// RequiredWith:  []string{"ldap_url", "ldap_base_dn", "ldap_uid", "ldap_verify_cert"},
			},
			"ldap_base_dn": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"oidc_name", "oidc_endpoint", "oidc_client_id", "oidc_client_secret", "oidc_groups_claim", "oidc_scope", "oidc_verify_cert", "oidc_auto_onboard", "oidc_user_claim"},
				RequiredWith:  []string{"ldap_url", "ldap_base_dn", "ldap_uid", "ldap_verify_cert"},
			},
			"ldap_uid": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"oidc_name", "oidc_endpoint", "oidc_client_id", "oidc_client_secret", "oidc_groups_claim", "oidc_scope", "oidc_verify_cert", "oidc_auto_onboard", "oidc_user_claim"},
				RequiredWith:  []string{"ldap_url", "ldap_base_dn", "ldap_uid", "ldap_verify_cert"},
			},

			"ldap_verify_cert": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"oidc_name", "oidc_endpoint", "oidc_client_id", "oidc_client_secret", "oidc_groups_claim", "oidc_scope", "oidc_verify_cert", "oidc_auto_onboard", "oidc_user_claim"},
				RequiredWith:  []string{"ldap_url", "ldap_base_dn", "ldap_uid", "ldap_verify_cert"},
			},
			"ldap_search_dn": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"oidc_name", "oidc_endpoint", "oidc_client_id", "oidc_client_secret", "oidc_groups_claim", "oidc_scope", "oidc_verify_cert", "oidc_auto_onboard", "oidc_user_claim"},
				RequiredWith:  []string{"ldap_url", "ldap_base_dn", "ldap_uid", "ldap_verify_cert"},
			},
			"ldap_search_password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"oidc_name", "oidc_endpoint", "oidc_client_id", "oidc_client_secret", "oidc_groups_claim", "oidc_scope", "oidc_verify_cert", "oidc_auto_onboard", "oidc_user_claim"},
				RequiredWith:  []string{"ldap_url", "ldap_base_dn", "ldap_uid", "ldap_verify_cert"},
			},
			"ldap_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ldap_group_uid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ldap_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ldap_group_base_dn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ldap_group_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ldap_group_gid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ldap_group_admin_dn": {
				Type:     schema.TypeString,
				Optional: true},
			"ldap_group_membership": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ldap_group_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Create: resourceConfigAuthCreate,
		Read:   resourceConfigAuthRead,
		Update: resourceConfigAuthCreate,
		Delete: resourceConfigAuthDelete,
	}
}

func resourceConfigAuthCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.GetConfigAuth(d)

	_, _, err := apiClient.SendRequest("PUT", models.PathConfig, body, 200)
	if err != nil {
		return err
	}

	return resourceConfigAuthRead(d, m)
}

func resourceConfigAuthRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, err := apiClient.SendRequest("GET", models.PathConfig, nil, 200)
	if err != nil {
		return err
	}

	var jsonData models.ConfigBodyResponse
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return err
	}

	auth := jsonData.AuthMode.Value
	d.Set("auth_mode", auth)

	switch auth {
	case "ldap_auth", "ldap":
		d.Set("ldap_url", jsonData.LdapURL.Value)
		d.Set("ldap_base_dn", jsonData.LdapBaseDn.Value)
		d.Set("ldap_uid", jsonData.LdapUID.Value)
		d.Set("ldap_verify_cert", jsonData.VerifyRemoteCert.Value)
		d.Set("ldap_search_dn", jsonData.LdapSearchDn.Value)
		d.Set("ldap_scope", jsonData.LdapScope.Value)
		d.Set("ldap_group_base_dn", jsonData.LdapGroupBaseDn.Value)
		d.Set("ldap_group_filter", jsonData.LdapGroupSearchFilter)
		d.Set("ldap_group_gid", jsonData.LdapGroupAttributeName.Value)
		d.Set("ldap_group_admin_dn", jsonData.LdapGroupAdminDn)
		d.Set("ldap_group_membership", jsonData.LdapGroupMembershipAttribute)
		d.Set("ldap_group_scope", jsonData.LdapGroupSearchScope)
	}

	d.SetId("/configurations")
	return nil
}

// func resourceConfigAuthUpdate(d *schema.ResourceData, m interface{}) error {
// 	apiClient := m.(*client.Client)
// 	body := client.GetConfigAuth(d)

// 	_, _, err := apiClient.SendRequest("PUT", models.PathConfig, body, 200)
// 	if err != nil {
// 		return err
// 	}

// 	return resourceConfigAuthRead(d, m)
// }

func resourceConfigAuthDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
