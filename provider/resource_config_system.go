package provider

import (
	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConfigSystem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_creation_restriction": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "adminonly",
			},
			"read_only": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "false",
			},
			"robot_token_expiration": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
		},
		Create: resourceConfigSystemCreate,
		Read:   resourceConfigSystemRead,
		Update: resourceConfigSystemUpdate,
		Delete: resourceConfigSystemDelete,
	}
}

func resourceConfigSystemCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.GetConfigBody(d)

	_, _, err := apiClient.SendRequest("PUT", models.PathConfig, body, 200)
	if err != nil {
		return err
	}

	return resourceConfigSystemRead(d, m)
}

func resourceConfigSystemRead(d *schema.ResourceData, m interface{}) error {
	d.SetId("configuration/system")
	return nil
}

func resourceConfigSystemUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.GetConfigBody(d)

	_, _, err := apiClient.SendRequest("PUT", models.PathConfig, body, 200)
	if err != nil {
		return err
	}

	return resourceConfigSystemRead(d, m)
}

func resourceConfigSystemDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
