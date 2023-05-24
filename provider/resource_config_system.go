package provider

import (
	"encoding/json"
	"fmt"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
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
				Default:  "robot$",
			},
			"scanner_skip_update_pulltime": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

	_, _, _, err := apiClient.SendRequest("PUT", models.PathConfig, body, 200)
	if err != nil {
		return err
	}

	return resourceConfigSystemRead(d, m)
}

func resourceConfigSystemRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest("GET", models.PathConfig, nil, 200)
	if respCode == 404 && err != nil {
		d.SetId("")
		return fmt.Errorf("Error getting system configuration %s", err)
	}

	var jsonData models.ConfigBodyResponse
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("Error getting system configuration %s", err)
	}

	d.Set("project_creation_restriction", jsonData.ProjectCreationRestriction.Value)
	d.Set("read_only", jsonData.ReadOnly.Value)
	d.Set("robot_token_expiration", jsonData.RobotTokenDuration.Value)
	d.Set("robot_name_prefix", jsonData.RobotNamePrefix.Value)
	d.Set("scanner_skip_update_pulltime", jsonData.ScannerSkipUpdatePulltime.Value)

	d.SetId("configuration/system")
	return nil
}

func resourceConfigSystemDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
