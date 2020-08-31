package provider

import (
	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Create: resourceConfigAuthCreate,
		Read:   resourceConfigAuthRead,
		Update: resourceConfigAuthUpdate,
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

	d.SetId("configuration/auth")
	return nil
}

func resourceConfigAuthUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.GetConfigAuth(d)

	_, _, err := apiClient.SendRequest("PUT", models.PathConfig, body, 200)
	if err != nil {
		return err
	}

	return resourceConfigAuthRead(d, m)
}

func resourceConfigAuthDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
