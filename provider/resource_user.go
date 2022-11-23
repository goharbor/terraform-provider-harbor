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

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"full_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"admin": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	body := client.UserBody(d)

	_, header, _, err := apiClient.SendRequest(ctx, "POST", models.PathUsers, &body, 201)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := client.GetID(header)
	if err != nil {
		return nil
	}

	d.SetId(id)
	return resourceUserRead(ctx, d, m)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	resp, _, respCode, err := apiClient.SendRequest(ctx, "GET", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] User %q was not found - removing from state!", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("making read request on user %s : %+v", d.Id(), err)
	}

	var jsonData models.UserBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("username", jsonData.Username)
	d.Set("full_name", jsonData.Realname)
	d.Set("email", jsonData.Email)
	d.Set("admin", jsonData.SysadminFlag)
	d.Set("comment", jsonData.Comment)

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	body := client.UserBody(d)
	_, _, _, err := apiClient.SendRequest(ctx, "PUT", d.Id(), body, 200)
	if err != nil {
		return diag.FromErr(err)
	}

	_, _, _, err = apiClient.SendRequest(ctx, "PUT", d.Id()+"/sysadmin", body, 200)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("password") == true {
		_, _, _, err = apiClient.SendRequest(ctx, "PUT", d.Id()+"/password", body, 200)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	_, _, respCode, err := apiClient.SendRequest(ctx, "DELETE", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] User %q was not found - already deleted!", d.Id())
			return nil
		}
		return diag.Errorf("making delete request on user %s : %+v", d.Id(), err)
	}
	return nil
}
