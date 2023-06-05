package provider

import (
	"encoding/json"
	"fmt"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConfigSecurity() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cve_allowlist": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"expires_at": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Create: resourceConfigSecurityCreate,
		Read:   resourceConfigSecurityRead,
		Update: resourceConfigSecurityUpdate,
		Delete: resourceConfigSecurityDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceConfigSecurityCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.SystemCVEAllowListBody(d)

	_, _, _, err := apiClient.SendRequest("PUT", models.PathSystemCVEAllowList, body, 200)
	if err != nil {
		return err
	}

	return resourceConfigSecurityRead(d, m)
}

func resourceConfigSecurityRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest("GET", models.PathSystemCVEAllowList, nil, 200)
	if respCode == 404 && err != nil {
		d.SetId("")
		return fmt.Errorf("resource not found %s", models.PathSystemCVEAllowList)
	}

	var jsonData models.SystemCveAllowListBodyPost
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("resource not found %s", models.PathSystemCVEAllowList)
	}

	// Convert the list of SystemCveAllowlistItems to a list of strings
	allowlistItems := make([]string, len(jsonData.Items))
	for i, item := range jsonData.Items {
		allowlistItems[i] = item.CveID
	}

	d.SetId(models.PathSystemCVEAllowList)
	d.Set("update_time", jsonData.UpdateTime)
	d.Set("expires_at", jsonData.ExpiresAt)
	d.Set("cve_allowlist", allowlistItems)
	d.Set("creation_time", jsonData.CreationTime)

	return nil
}

func resourceConfigSecurityUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.SystemCVEAllowListBody(d)

	_, _, _, err := apiClient.SendRequest("PUT", models.PathSystemCVEAllowList, body, 200)
	if err != nil {
		return err
	}

	return resourceConfigSecurityRead(d, m)
}

func resourceConfigSecurityDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	d.Set("expires_at", nil)
	d.Set("cve_allowlist", []string{})
	body := client.SystemCVEAllowListBody(d)

	_, _, _, err := apiClient.SendRequest("PUT", models.PathSystemCVEAllowList, body, 200)
	if err != nil {
		return err
	}

	return resourceConfigSecurityRead(d, m)
}
