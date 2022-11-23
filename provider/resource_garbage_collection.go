package provider

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"

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
		},
		CreateContext: resourceGCCreateUpdate,
		ReadContext:   resourceGCRead,
		UpdateContext: resourceGCCreateUpdate,
		DeleteContext: resourceGCDelete,
	}
}

func resourceGCCreateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	err := apiClient.SetSchedule(ctx, d, "gc")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("/system/gc/schedule")
	return resourceGCRead(ctx, d, m)
}

func resourceGCRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest(ctx, "GET", models.PathGC, nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Garbage collection %q was not found - removing from state!", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("making Read request on garbage collection %s : %+v", d.Id(), err)
	}

	var jsonData models.SystemBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return diag.FromErr(err)
	}

	jobParameters := jsonData.JobParameters

	var jsonJobParameters models.JobParameters
	err = json.Unmarshal([]byte(jobParameters), &jsonJobParameters)
	if err != nil {
		return diag.FromErr(err)
	}

	if jsonData.Schedule.Type == "Custom" {
		d.Set("schedule", jsonData.Schedule.Cron)
	} else {
		d.Set("schedule", jsonData.Schedule.Type)
	}
	d.Set("delete_untagged", jsonJobParameters.DeleteUntagged)
	return nil
}

func resourceGCDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	body := models.SystemBody{}
	body.Schedule.Cron = ""
	body.Schedule.Type = "None"
	body.Parameters.DeleteUntagged = false

	_, _, _, err := apiClient.SendRequest(ctx, "PUT", models.PathGC, body, 200)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
