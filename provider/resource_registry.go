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
		CreateContext: resourceRegistryCreate,
		ReadContext:   resourceRegistryRead,
		UpdateContext: resourceRegistryUpdate,
		DeleteContext: resourceRegistryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceRegistryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	body := client.GetRegistryBody(d)

	_, headers, _, err := apiClient.SendRequest(ctx, "POST", models.PathRegistries, body, 201)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := client.GetID(headers)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return resourceRegistryRead(ctx, d, m)
}

func resourceRegistryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest(ctx, "GET", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Registry %q was not found - removing from state!", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("making read request on registry %s : %+v", d.Id(), err)
	}

	var jsonData models.RegistryBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return diag.FromErr(err)
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

func resourceRegistryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	body := client.GetRegistryBody(d)

	_, _, _, err := apiClient.SendRequest(ctx, "PUT", d.Id(), body, 200)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRegistryRead(ctx, d, m)
}

func resourceRegistryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	_, _, respCode, err := apiClient.SendRequest(ctx, "DELETE", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Registry %q was not found - already deleted!", d.Id())
			return nil
		}
		return diag.Errorf("making delete request on registry %s : %+v", d.Id(), err)
	}
	return nil
}
