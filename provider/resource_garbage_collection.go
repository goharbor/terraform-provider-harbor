package provider

import (
	"encoding/json"
	"fmt"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"workers": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				ValidateFunc: func(i interface{}, k string) (ws []string, errors []error) {
					value := i.(int)
					if value < 1 || value > 5 {
						errors = append(errors, fmt.Errorf("GC workers must be between 1 and 5"))
					}
					return
				},
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

	resp, _, respCode, err := apiClient.SendRequest("GET", models.PathGC, nil, 200)
	if respCode == 404 && err != nil {
		d.SetId("")
		return fmt.Errorf("Resource not found %s", d.Id())
	}

	var jsonData models.SystemBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return err
	}

	jobParameters := jsonData.JobParameters

	var jsonJobParameters models.JobParameters
	err = json.Unmarshal([]byte(jobParameters), &jsonJobParameters)
	if err != nil {
		fmt.Println(err)
	}

	if jsonData.Schedule.Type == "Custom" {
		d.Set("schedule", jsonData.Schedule.Cron)
	} else {
		d.Set("schedule", jsonData.Schedule.Type)
	}
	d.Set("delete_untagged", jsonJobParameters.DeleteUntagged)
	d.Set("workers", jsonJobParameters.Workers)
	return nil
}

func resourceGCDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := models.SystemBody{}
	body.Schedule.Cron = ""
	body.Schedule.Type = "None"
	body.Parameters.DeleteUntagged = false

	_, _, _, err := apiClient.SendRequest("PUT", models.PathGC, body, 200)
	if err != nil {
		return err
	}
	return nil
}
