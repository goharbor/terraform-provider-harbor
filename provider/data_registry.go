package provider

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataRegistry() *schema.Resource {
	return &schema.Resource{
		Read: dataRegistryRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"registry_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"insecure": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataRegistryRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	name := d.Get("name").(string)
	registryPath := models.PathRegistries + "?name=" + name
	resp, _, _, err := apiClient.SendRequest("GET", registryPath, nil, 200)
	if err != nil {
		return err
	}

	var jsonData []models.RegistryBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("Unable to find the registry named: %s", name)
	}

	for _, v := range jsonData {
		if v.Name == name {
			id := models.PathProjects + "/" + strconv.Itoa(v.ID)

			d.SetId(id)
			d.Set("registry_id", v.ID)
			d.Set("name", v.Name)
			d.Set("type", v.Type)
			d.Set("description", v.Description)
			d.Set("url", v.URL)
			d.Set("insecure", v.Insecure)
			d.Set("status", v.Status)
		}
	}

	return nil
}
