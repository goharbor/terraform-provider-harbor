package provider

import (
	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var pathConfig = "/configurations"

type system struct {
	ProjectCreationRestriction string `json:"project_creation_restriction"`
	ReadOnly                   string `json:"read_only,omitempty"`
	RobotTokenDuration         int    `json:"robot_token_duration,omitempty"`

	// EmailVerifyCert string `json:"email_verify_cert,omitempty"`
}

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

func days2mins(days int) int {
	mins := days * 1440
	return mins
}

func resourceConfigSystemCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := system{
		ProjectCreationRestriction: d.Get("project_creation_restriction").(string),
		ReadOnly:                   d.Get("read_only").(string),
		RobotTokenDuration:         days2mins(d.Get("robot_token_expiration").(int)),
	}

	_, err := apiClient.SendRequest("PUT", pathConfig, body, 200)
	if err != nil {
		return err
	}

	d.SetId(randomString(15))
	// return resourceConfigSystemRead(d, m)
	return nil
}

func resourceConfigSystemRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceConfigSystemUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := system{
		ProjectCreationRestriction: d.Get("project_creation_restriction").(string),
		ReadOnly:                   d.Get("read_only").(string),
		RobotTokenDuration:         days2mins(d.Get("robot_token_expiration").(int)),
	}

	_, err := apiClient.SendRequest("PUT", pathConfig, body, 200)
	if err != nil {
		return err
	}

	return resourceConfigSystemRead(d, m)
}

func resourceConfigSystemDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
