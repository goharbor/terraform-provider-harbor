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

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
		CreateContext: resourceGroupCreateUpdate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupCreateUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceGroupCreateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	body := client.GroupBody(d)

	_, header, _, err := apiClient.SendRequest(ctx, "POST", models.PathGroups, &body, 201)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := client.GetID(header)
	if err != nil {
		return nil
	}

	d.SetId(id)
	return resourceGroupRead(ctx, d, m)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	resp, _, respCode, err := apiClient.SendRequest(ctx, "GET", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Resource group %q was not found - removing from state!", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("making Read request on resource group %s : %+v", d.Id(), err)
	}

	var jsonData models.GroupBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("group_name", jsonData.Groupname)
	d.Set("group_type", jsonData.GroupType)

	return nil
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	_, _, respCode, err := apiClient.SendRequest(ctx, "DELETE", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Resource group %q was not found - already deleted!", d.Id())
			return nil
		}
		return diag.Errorf("making delete request on resource group %s : %+v", d.Id(), err)
	}

	return nil
}
