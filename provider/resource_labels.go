package provider

import (
	"encoding/json"
	"fmt"

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
		Create: resourceLabelCreate,
		Read:   resourceLabelRead,
		Update: resourceLabelUpdate,
		Delete: resourceLabelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceLabelCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.LabelsBody(d)

	_, headers, _, err := apiClient.SendRequest("POST", models.PathLabel, body, 201)
	if err != nil {
		return err
	}

	id, err := client.GetID(headers)
	d.SetId(id)
	return resourceLabelRead(d, m)
}

func resourceLabelRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	if respCode == 404 && err != nil {
		d.SetId("")
		return nil
	} else if err != nil {
		return fmt.Errorf("resource not found %s", d.Id())
	}

	var jsonData models.Labels
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("Resource not found %s", d.Id())
	}

	d.Set("name", jsonData.Name)
	d.Set("description", jsonData.Description)
	d.Set("color", jsonData.Color)
	d.Set("scope", jsonData.Scope)

	return nil
}

func resourceLabelUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.LabelsBody(d)

	_, _, _, err := apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		return err
	}

	return resourceLabelRead(d, m)
}

func resourceLabelDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	_, _, respCode, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if respCode != 404 && err != nil { // We can't delete something that doesn't exist. Hence the 404-check
		return err
	}
	return nil
}
