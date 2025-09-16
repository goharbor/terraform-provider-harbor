package provider

import (
	"encoding/json"
	"fmt"

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
				Default:  "",
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
			"payload_format": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Default",
			},
		},
		Create: resourceProjectWebhookCreate,
		Read:   resourceProjectWebhookRead,
		Update: resourceProjectWebhookUpdate,
		Delete: resourceProjectWebhookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceProjectWebhookCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.ProjectWebhookBody(d)

	url := d.Get("project_id").(string) + "/webhook/policies"
	_, headers, _, err := apiClient.SendRequest("POST", url, body, 201)
	if err != nil {
		return err
	}

	id, _ := client.GetID(headers)
	d.SetId(id)
	return resourceProjectWebhookRead(d, m)
}

func resourceProjectWebhookRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	if respCode == 404 && err != nil {
		d.SetId("")
		return fmt.Errorf("resource not found %s", d.Id())
	} else if err != nil {
		return err
	}

	var jsonData models.ProjectWebhook
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("resource not found %s", d.Id())
	}

	d.Set("name", jsonData.Name)
	d.Set("description", jsonData.Description)
	d.Set("enabled", jsonData.Enabled)
	d.Set("notify_type", jsonData.Targets[0].Type)
	d.Set("address", jsonData.Targets[0].Address)
	d.Set("auth_header", jsonData.Targets[0].AuthHeader)
	d.Set("skip_cert_verify", jsonData.Targets[0].SkipCertVerify)
	d.Set("payload_format", jsonData.Targets[0].PayloadFormat)

	return nil
}

func resourceProjectWebhookUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.ProjectWebhookBody(d)

	_, _, _, err := apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		return err
	}

	return resourceProjectWebhookRead(d, m)
}

func resourceProjectWebhookDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	_, _, _, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if err != nil {
		return err
	}
	return nil
}
