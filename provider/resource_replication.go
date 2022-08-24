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

func resourceReplication() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"deletion": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"schedule": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "manual",
			},
			"registry_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"replication_policy_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"dest_namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"override": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"filters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tag": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"labels": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"decoration": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
		CreateContext: resourceReplicationCreate,
		ReadContext:   resourceReplicationRead,
		UpdateContext: resourceReplicationUpdate,
		DeleteContext: resourceReplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceReplicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	body := client.GetReplicationBody(d)

	_, headers, _, err := apiClient.SendRequest(ctx, "POST", models.PathReplication, body, 201)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := client.GetID(headers)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return resourceReplicationRead(ctx, d, m)
}

func resourceReplicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest(ctx, "GET", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Replication %q was not found - removing from state!", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("making Read request on replication m %s : %+v", d.Id(), err)
	}

	var jsonData models.RegistryBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return diag.FromErr(err)
	}

	var jsonDataReplication models.ReplicationBody
	err = json.Unmarshal([]byte(resp), &jsonDataReplication)
	if err != nil {
		return diag.FromErr(err)
	}

	destRegistryID := jsonDataReplication.DestRegistry.ID

	if destRegistryID == 0 {
		d.Set("action", "pull")
		d.Set("registry_id", jsonDataReplication.SrcRegistry.ID)

	} else {
		d.Set("action", "push")
		d.Set("registry_id", destRegistryID)
	}

	switch jsonDataReplication.Trigger.Type {
	case "scheduled":
		d.Set("schedule", jsonDataReplication.Trigger.TriggerSettings.Cron)
	case "event_based":
		d.Set("schedule", "event_based")
	default:
		d.Set("schedule", "manual")
	}

	d.Set("replication_policy_id", jsonDataReplication.ID)
	d.Set("enabled", jsonDataReplication.Enabled)
	d.Set("name", jsonDataReplication.Name)
	d.Set("deletion", jsonDataReplication.Deletion)
	d.Set("override", jsonDataReplication.Override)

	return nil
}

func resourceReplicationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	body := client.GetReplicationBody(d)

	_, _, _, err := apiClient.SendRequest(ctx, "PUT", d.Id(), body, 200)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceReplicationRead(ctx, d, m)
}

func resourceReplicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	_, _, respCode, err := apiClient.SendRequest(ctx, "DELETE", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Replication %q was not found - already deleted!", d.Id())
			return nil
		}
		return diag.Errorf("making delete request on replication %s : %+v", d.Id(), err)
	}
	return nil
}
