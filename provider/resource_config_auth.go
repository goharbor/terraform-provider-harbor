package provider

import (
	"context"
	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceConfigAuth() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"auth_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"oidc_name": {
				Type:          schema.TypeString,
				Optional:      true,
				RequiredWith:  oidcRequiredWith(),
				ConflictsWith: oidcConflictsWith(),
			},
			"oidc_endpoint": {
				Type:          schema.TypeString,
				Optional:      true,
				RequiredWith:  oidcRequiredWith(),
				ConflictsWith: oidcConflictsWith(),
			},
			"oidc_client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				RequiredWith:  oidcRequiredWith(),
				ConflictsWith: oidcConflictsWith(),
			},
			"oidc_client_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				RequiredWith:  oidcRequiredWith(),
				ConflictsWith: oidcConflictsWith(),
			},
			"oidc_groups_claim": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: oidcConflictsWith(),
			},
			"oidc_admin_group": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: oidcConflictsWith(),
			},
			"oidc_scope": {
				Type:          schema.TypeString,
				Optional:      true,
				RequiredWith:  oidcRequiredWith(),
				ConflictsWith: oidcConflictsWith(),
			},
			"oidc_verify_cert": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: oidcConflictsWith(),
			},
			"oidc_auto_onboard": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: oidcConflictsWith(),
			},
			"oidc_user_claim": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: oidcConflictsWith(),
			},
			"ldap_url": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: ldapConflictsWith(),
				RequiredWith:  ldapRequiredWith(),
			},
			"ldap_base_dn": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: ldapConflictsWith(),
				RequiredWith:  ldapRequiredWith(),
			},
			"ldap_uid": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: ldapConflictsWith(),
				RequiredWith:  ldapRequiredWith(),
			},

			"ldap_verify_cert": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: ldapConflictsWith(),
				RequiredWith:  ldapRequiredWith(),
			},
			"ldap_search_dn": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: ldapConflictsWith(),
				RequiredWith:  ldapRequiredWith(),
			},
			"ldap_search_password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: ldapConflictsWith(),
				RequiredWith:  ldapRequiredWith(),
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
		CreateContext: resourceConfigAuthCreateUpdate,
		ReadContext:   resourceConfigAuthRead,
		UpdateContext: resourceConfigAuthCreateUpdate,
		DeleteContext: resourceConfigAuthDelete,
	}
}

func resourceConfigAuthCreateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	body := client.GetConfigAuth(d)

	_, _, _, err := apiClient.SendRequest(ctx, "PUT", models.PathConfig, body, 200)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceConfigAuthRead(ctx, d, m)
}

func resourceConfigAuthRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest(ctx, "GET", models.PathConfig, nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Config Auth %q was not found - removing from state!", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("making Read request on Config Auth %s : %+v", d.Id(), err)
	}

	err = client.SetAuthValues(d, resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(models.PathConfig)
	return nil
}

func resourceConfigAuthDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return nil
}

func oidcConflictsWith() []string {
	return []string{"ldap_url",
		"ldap_base_dn",
		"ldap_uid",
		"ldap_verify_cert",
		"ldap_search_dn",
		"ldap_search_password",
		"ldap_filter",
		"ldap_group_uid",
		"ldap_scope",
		"ldap_group_base_dn",
		"ldap_group_filter",
		"ldap_group_gid",
		"ldap_group_admin_dn",
		"ldap_group_membership",
		"ldap_group_scope"}
}

func oidcRequiredWith() []string {
	return []string{"oidc_name", "oidc_endpoint", "oidc_client_id", "oidc_client_secret", "oidc_scope"}
}

func ldapConflictsWith() []string {
	return []string{"oidc_name", "oidc_endpoint", "oidc_client_id", "oidc_client_secret", "oidc_groups_claim", "oidc_scope", "oidc_verify_cert", "oidc_auto_onboard", "oidc_user_claim"}
}

func ldapRequiredWith() []string {
	return []string{"ldap_url", "ldap_base_dn", "ldap_uid", "ldap_verify_cert"}
}
