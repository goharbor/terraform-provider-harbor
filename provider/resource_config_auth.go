package provider

import (
	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type auth struct {
	AuthMode         string `json:"auth_mode"`
	OidcName         string `json:"oidc_name,omitempty"`
	OidcEndpoint     string `json:"oidc_endpoint,omitempty"`
	OidcClientID     string `json:"oidc_client_id,omitempty"`
	OidcClientSecret string `json:"oidc_client_secret,omitempty"`
	OidcScope        string `json:"oidc_scope,omitempty"`
	OidcVerifyCert   string `json:"oidc_verify_cert,omitempty"`
	OidcGroupsClaim  string `json:"oidc_groups_claim,omitempty"`
}

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
	apiClient, body := newAPIClient(d, m)
	_, err := apiClient.SendRequest("PUT", pathConfig, body, 200)
	if err != nil {
		return err
	}

	d.SetId(randomString(15))
	// return resourceConfigAuthRead(d, m)
	return nil
}

func resourceConfigAuthRead(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceConfigAuthUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient, body := newAPIClient(d, m)

	_, err := apiClient.SendRequest("PUT", pathConfig, body, 200)
	if err != nil {
		return err
	}

	return resourceConfigAuthRead(d, m)
}

func resourceConfigAuthDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func newAPIClient(d *schema.ResourceData, m interface{}) (*client.Client, auth) {
	apiClient := m.(*client.Client)

	body := auth{
		AuthMode:         d.Get("auth_mode").(string),
		OidcName:         d.Get("oidc_name").(string),
		OidcEndpoint:     d.Get("oidc_endpoint").(string),
		OidcClientID:     d.Get("oidc_client_id").(string),
		OidcClientSecret: d.Get("oidc_client_secret").(string),
		OidcGroupsClaim:  d.Get("oidc_groups_claim").(string),
		OidcScope:        d.Get("oidc_scope").(string),
		OidcVerifyCert:   d.Get("oidc_verify_cert").(string),
	}
	return apiClient, body
}
