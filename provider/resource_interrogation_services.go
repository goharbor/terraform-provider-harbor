package provider

import (
	"context"
	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVuln() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vulnerability_scan_policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_scanner": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		CreateContext: resourceVulnCreateUpdate,
		ReadContext:   resourceVulnRead,
		UpdateContext: resourceVulnCreateUpdate,
		DeleteContext: resourceVulnDelete,
	}
}

func resourceVulnCreateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	err := apiClient.SetSchedule(ctx, d, "vuln")
	if err != nil {
		return diag.FromErr(err)
	}

	scanner := d.Get("default_scanner").(string)
	if scanner != "" {
		apiClient.SetDefaultScanner(ctx, scanner)
	}

	return resourceVulnRead(ctx, d, m)
}

func resourceVulnRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("/system/scanAll/schedule")
	return nil
}

func resourceVulnDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	body := models.SystemBody{}
	body.Schedule.Cron = ""
	body.Schedule.Type = "None"

	_, _, _, err := apiClient.SendRequest(ctx, "PUT", models.PathVuln, body, 200)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
