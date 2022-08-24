package provider

import (
	"context"
	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConfigSystem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_creation_restriction": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "adminonly",
			},
			"read_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"robot_token_expiration": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"robot_name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		CreateContext: resourceConfigSystemCreateUpdate,
		ReadContext:   resourceConfigSystemRead,
		UpdateContext: resourceConfigSystemCreateUpdate,
		DeleteContext: resourceConfigSystemDelete,
	}
}

func resourceConfigSystemCreateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	body := client.GetConfigSystem(d)

	_, _, _, err := apiClient.SendRequest(ctx, "PUT", models.PathConfig, body, 200)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceConfigSystemRead(ctx, d, m)
}

func resourceConfigSystemRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("configuration/system")
	return nil
}

func resourceConfigSystemDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
