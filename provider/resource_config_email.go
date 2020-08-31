package provider

import (
	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConfigEmail() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"email_host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email_port": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email_username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"email_from": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email_ssl": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "false",
			},
			// "email_verify_cert": {
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// },
		},
		Create: resourceConfigEmailCreate,
		Read:   resourceConfigEmailRead,
		Update: resourceConfigEmailUpdate,
		Delete: resourceConfigEmailDelete,
	}
}

func resourceConfigEmailCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.GetConfigEmail(d)

	_, _, err := apiClient.SendRequest("PUT", models.PathConfig, body, 200)
	if err != nil {
		return err
	}

	return resourceConfigEmailRead(d, m)
}

func resourceConfigEmailRead(d *schema.ResourceData, m interface{}) error {
	d.SetId("configuration/email")
	return nil
}

func resourceConfigEmailUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.GetConfigEmail(d)

	_, _, err := apiClient.SendRequest("PUT", models.PathConfig, body, 200)
	if err != nil {
		return err
	}

	return resourceConfigEmailRead(d, m)
}

func resourceConfigEmailDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
