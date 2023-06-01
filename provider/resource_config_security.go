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
		return fmt.Errorf("resource not found %s", d.Id())
	}


	expires_at, expires_at_true := d.GetOk("expires_at")

	if expires_at_true {
		d.Set("expires_at", expires_at)
	}

	d.SetId(models.PathSystemCVEAllowList)
	d.Set("update_time", jsonData.UpdateTime)
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

func resourceConfigSecurityDelete(d *schema.ResourceData, m interface{}) error { // Harbor doesn't really have aby way to delete this resource yet...
	apiClient := m.(*client.Client)
	body := client.SystemCVEAllowListBody(d)

	_, _, _, err := apiClient.SendRequest("PUT", models.PathSystemCVEAllowList, body, 200)
	if err != nil {
		return err
	}

	return resourceConfigSecurityRead(d, m)
}
