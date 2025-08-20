package provider

import (
	"encoding/json"
	"fmt"

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
	apiClient := m.(*client.Client)
	resp, _, respCode, err := apiClient.SendRequest("GET", d.Id(), nil, 200)

	if respCode == 404 && err != nil {
		d.SetId("")
		return fmt.Errorf("resource not found %s", d.Id())
	} else if err != nil {
		return err
	}

	var jsonData models.InterogationsBodyResponse
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("resource not found %s", d.Id())
	}

	vulnerability_scan_policy := jsonData.Schedule.Type
	if vulnerability_scan_policy == "Custom" {
		vulnerability_scan_policy = jsonData.Schedule.Cron
	}

	d.Set("vulnerability_scan_policy", vulnerability_scan_policy)

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
