package provider

import (
	"encoding/json"
	"fmt"

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
			"dest_namespace_replace": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
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
			"speed": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
		Create: resourceReplicationCreate,
		Read:   resourceReplicationRead,
		Update: resourceReplicationUpdate,
		Delete: resourceReplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceReplicationCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.GetReplicationBody(d)

	_, headers, _, err := apiClient.SendRequest("POST", models.PathReplication, body, 201)
	if err != nil {
		return err
	}

	id, err := client.GetID(headers)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceReplicationRead(d, m)
}

func resourceReplicationRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	if respCode == 404 && err != nil {
		d.SetId("")
		return fmt.Errorf("Resource not found %s", d.Id())
	}
	var jsonData models.RegistryBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("Resource not found %s", d.Id())
	}

	var jsonDataReplication models.ReplicationBody
	err = json.Unmarshal([]byte(resp), &jsonDataReplication)
	if err != nil {
		return fmt.Errorf("Resource not found %s", d.Id())
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

func resourceReplicationUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.GetReplicationBody(d)

	_, _, _, err := apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		return err
	}

	return resourceReplicationRead(d, m)
}

func resourceReplicationDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	_, _, _, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if err != nil {
		return err
	}
	return nil
}
