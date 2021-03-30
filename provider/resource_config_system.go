package provider

import (
	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"robot_token_expiration": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"robot_name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Create: resourceConfigSystemCreate,
		Read:   resourceConfigSystemRead,
		Update: resourceConfigSystemCreate,
		Delete: resourceConfigSystemDelete,
	}
}

func resourceConfigSystemCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.GetConfigSystem(d)

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

// func resourceConfigSystemUpdate(d *schema.ResourceData, m interface{}) error {
// 	apiClient := m.(*client.Client)

// 	body := client.GetConfigSystem(d)

// 	_, _, err := apiClient.SendRequest("PUT", models.PathConfig, body, 200)
// 	if err != nil {
// 		return err
// 	}

// 	return resourceConfigSystemRead(d, m)
// }

func resourceConfigSystemDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
