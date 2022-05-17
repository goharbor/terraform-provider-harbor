package provider

import (
	"encoding/json"
	"fmt"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRegistry() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoint_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_secret": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"insecure": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"registry_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		Create: resourceRegistryCreate,
		Read:   resourceRegistryRead,
		Update: resourceRegistryUpdate,
		Delete: resourceRegistryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceRegistryCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.GetRegistryBody(d)

	_, headers, _, err := apiClient.SendRequest("POST", models.PathRegistries, body, 201)
	if err != nil {
		return err
	}

	id, err := client.GetID(headers)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceRegistryRead(d, m)
}

func resourceRegistryRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, _, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	if err != nil {
		d.SetId("")
		return nil
	}

	var jsonData models.RegistryBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Resource not found %s", d.Id())
	}

	registryName, _ := client.GetRegistryType(jsonData.Type)

	d.Set("name", jsonData.Name)
	d.Set("description", jsonData.Description)
	d.Set("endpoint_url", jsonData.URL)
	d.Set("access_id", jsonData.Credential.AccessKey)
	d.Set("insecure", jsonData.Insecure)
	d.Set("status", jsonData.Status)
	d.Set("registry_id", jsonData.ID)
	d.Set("provider_name", registryName)

	return nil
}

func resourceRegistryUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.GetRegistryBody(d)

	_, _, _, err := apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		return err
	}

	return resourceRegistryRead(d, m)
}

func resourceRegistryDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	_, _, _, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if err != nil {
		return err
	}
	return nil
}
