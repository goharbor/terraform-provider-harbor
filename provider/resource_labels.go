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

func resourceLabel() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"color": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "false",
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		CreateContext: resourceLabelCreate,
		ReadContext:   resourceLabelRead,
		UpdateContext: resourceLabelUpdate,
		DeleteContext: resourceLabelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceLabelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	body := client.LabelsBody(d)

	_, headers, _, err := apiClient.SendRequest(ctx, "POST", models.PathLabel, body, 201)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := client.GetID(headers)
	d.SetId(id)
	return resourceLabelRead(ctx, d, m)
}

func resourceLabelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest(ctx, "GET", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Label %q was not found - removing from state!", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("making Read request on label %s : %+v", d.Id(), err)
	}

	var jsonData models.Labels
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", jsonData.Name)
	d.Set("description", jsonData.Description)
	d.Set("color", jsonData.Color)
	d.Set("scope", jsonData.Scope)

	return nil
}

func resourceLabelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	body := client.LabelsBody(d)

	_, _, _, err := apiClient.SendRequest(ctx, "PUT", d.Id(), body, 200)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLabelRead(ctx, d, m)
}

func resourceLabelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	_, _, respCode, err := apiClient.SendRequest(ctx, "DELETE", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Label %q was not found - already deleted!", d.Id())
			return nil
		}
		return diag.Errorf("making delete request on label %s : %+v", d.Id(), err)
	}
	return nil
}
