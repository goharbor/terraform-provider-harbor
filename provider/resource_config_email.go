package provider

import (
	"context"
	"encoding/json"
	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceConfigEmail() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"email_host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  25,
			},
			"email_username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"email_from": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"email_insecure": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		CreateContext: resourceConfigEmailCreateUpdate,
		ReadContext:   resourceConfigEmailRead,
		UpdateContext: resourceConfigEmailCreateUpdate,
		DeleteContext: resourceConfigEmailDelete,
	}
}

func resourceConfigEmailCreateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	body := client.GetConfigEmail(d)

	_, _, _, err := apiClient.SendRequest(ctx, "PUT", models.PathConfig, body, 200)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceConfigEmailRead(ctx, d, m)
}

func resourceConfigEmailRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest(ctx, "GET", models.PathConfig, nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Config Email %q was not found - removing from state!", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("making Read request on Config Email %s : %+v", d.Id(), err)
	}

	var jsonData models.ConfigBodyResponse
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("/configurations")
	d.Set("email_host", jsonData.EmailHost.Value)
	d.Set("email_port", jsonData.EmailPort.Value)
	d.Set("email_username", jsonData.EmailUsername.Value)
	d.Set("email_from", jsonData.EmailFrom.Value)
	d.Set("email_ssl", jsonData.EmailSsl.Value)
	return nil
}

func resourceConfigEmailDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return nil
}
