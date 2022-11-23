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

func resourceProjectWebhook() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"events_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"notify_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auth_header": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"skip_cert_verify": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
		CreateContext: resourceProjectWebhookCreate,
		ReadContext:   resourceProjectWebhookRead,
		UpdateContext: resourceProjectWebhookUpdate,
		DeleteContext: resourceProjectWebhookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceProjectWebhookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	body := client.ProjectWebhookBody(d)

	url := d.Get("project_id").(string) + "/webhook/policies"
	_, headers, _, err := apiClient.SendRequest(ctx, "POST", url, body, 201)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := client.GetID(headers)
	d.SetId(id)
	return resourceProjectWebhookRead(ctx, d, m)
}

func resourceProjectWebhookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest(ctx, "GET", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Web hook %q was not found - removing from state!", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("making Read request on web hook %s : %+v", d.Id(), err)
	}

	var jsonData models.ProjectWebhook
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", jsonData.Name)
	d.Set("description", jsonData.Description)
	d.Set("enabled", jsonData.Enabled)
	d.Set("notify_type", jsonData.Targets[0].Type)
	d.Set("address", jsonData.Targets[0].Address)
	d.Set("auth_header", jsonData.Targets[0].AuthHeader)
	d.Set("skip_cert_verify", jsonData.Targets[0].SkipCertVerify)

	return nil
}

func resourceProjectWebhookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	body := client.ProjectWebhookBody(d)

	_, _, _, err := apiClient.SendRequest(ctx, "PUT", d.Id(), body, 200)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceProjectWebhookRead(ctx, d, m)
}

func resourceProjectWebhookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	_, _, respCode, err := apiClient.SendRequest(ctx, "DELETE", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Project webhook %q was not found - already deleted!", d.Id())
			return nil
		}
		return diag.Errorf("making delete request on project webhook %s : %+v", d.Id(), err)
	}

	return nil
}
