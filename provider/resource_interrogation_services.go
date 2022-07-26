package provider

import (
	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
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
		Create: resourceVulnCreate,
		Read:   resourceVulnRead,
		Update: resourceVulnCreate,
		Delete: resourceVulnDelete,
	}
}

func resourceVulnCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	err := apiClient.SetSchedule(d, "vuln")
	if err != nil {
		return err
	}

	scanner := d.Get("default_scanner").(string)
	if scanner != "" {
		apiClient.SetDefaultScanner(scanner)
	}

	return resourceVulnRead(d, m)
}

func resourceVulnRead(d *schema.ResourceData, m interface{}) error {
	d.SetId("/system/scanAll/schedule")
	return nil
}

func resourceVulnDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := models.SystemBody{}
	body.Schedule.Cron = ""
	body.Schedule.Type = "None"

	_, _, _, err := apiClient.SendRequest("PUT", models.PathVuln, body, 200)
	if err != nil {
		return err
	}
	return nil
}
