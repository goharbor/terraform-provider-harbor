package provider

import (
	"encoding/json"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGC() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"schedule": {
				Type:     schema.TypeString,
				Required: true,
			},
			"delete_untagged": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		Create: resourceGCCreate,
		Read:   resourceGCRead,
		Update: resourceGCCreate,
		Delete: resourceGCDelete,
	}
}

func resourceGCCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	err := apiClient.SetSchedule(d, "gc")
	if err != nil {
		return err
	}
	d.SetId("/system/gc/schedule")
	return resourceGCRead(d, m)
}

func resourceGCRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, err := apiClient.SendRequest("GET", models.PathGC, nil, 200)
	if err != nil {
		return err
	}

	var jsonData models.SystemBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return err
	}

	d.Set("schedule", jsonData.Schedule.Type)
	return nil
}

func resourceGCDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := models.SystemBody{}
	body.Schedule.Cron = ""
	body.Schedule.Type = "None"
	body.Parameters.DeleteUntagged = false

	_, _, err := apiClient.SendRequest("PUT", models.PathGC, body, 200)
	if err != nil {
		return err
	}
	return nil
}
