package provider

import (
	"encoding/json"
	"fmt"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		Create: resourceConfigEmailCreate,
		Read:   resourceConfigEmailRead,
		Update: resourceConfigEmailCreate,
		Delete: resourceConfigEmailDelete,
	}
}

func resourceConfigEmailCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.GetConfigEmail(d)

	_, _, _, err := apiClient.SendRequest("PUT", models.PathConfig, body, 200)
	if err != nil {
		return err
	}

	return resourceConfigEmailRead(d, m)
}

func resourceConfigEmailRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest("GET", models.PathConfig, nil, 200)
	if respCode == 404 && err != nil {
		d.SetId("")
		return fmt.Errorf("Resource not found %s", d.Id())
	}

	var jsonData models.ConfigBodyResponse
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return err
	}

	d.SetId("/configurations")
	d.Set("email_host", jsonData.EmailHost.Value)
	d.Set("email_port", jsonData.EmailPort.Value)
	d.Set("email_username", jsonData.EmailUsername.Value)
	d.Set("email_from", jsonData.EmailFrom.Value)
	d.Set("email_ssl", jsonData.EmailSsl.Value)
	return nil
}

func resourceConfigEmailDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
