package provider

import (
	"encoding/json"
	"fmt"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePreheatInstance() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "dragonfly" && v != "kraken" {
						errs = append(errs, fmt.Errorf("%q must be either 'dragonfly' or 'kraken', got: %s", key, v))
					}
					return
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auth_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NONE",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "NONE" && v != "BASIC" && v != "OAUTH" {
						errs = append(errs, fmt.Errorf("%q must be either 'NONE' or 'BASIC' or 'OAUTH', got: %s", key, v))
					}
					return
				},
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
				Default:   "",
			},
			"token": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
				Default:   "",
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"insecure": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		Create: resourcePreheatInstanceCreate,
		Read:   resourcePreheatInstanceRead,
		Update: resourcePreheatInstanceUpdate,
		Delete: resourcePreheatInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourcePreheatInstanceCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.PreheatInstanceBody(d)

	_, headers, _, err := apiClient.SendRequest("POST", models.PathPreheatInstance, body, 201)
	if err != nil {
		return err
	}

	id, err := client.GetID(headers)
	d.SetId(id)
	return resourcePreheatInstanceRead(d, m)
}

func resourcePreheatInstanceRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	if respCode == 404 && err != nil {
		d.SetId("")
		return nil
	} else if err != nil {
		return fmt.Errorf("resource not found %s", d.Id())
	}

	var jsonData models.PreheatInstance
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("Resource not found %s", d.Id())
	}

	d.Set("name", jsonData.Name)
	d.Set("description", jsonData.Description)
	d.Set("vendor", jsonData.Vendor)
	d.Set("endpoint", jsonData.Endpoint)
	d.Set("auth_mode", jsonData.AuthMode)
	d.Set("enabled", jsonData.Enabled)
	d.Set("default", jsonData.Default)
	d.Set("insecure", jsonData.Insecure)
	d.Set("username", jsonData.AuthInfo.Username)
	d.Set("password", jsonData.AuthInfo.Password)
	d.Set("token", jsonData.AuthInfo.Token)

	return nil
}

func resourcePreheatInstanceUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.PreheatInstanceBody(d)

	_, _, _, err := apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		return err
	}

	return resourcePreheatInstanceRead(d, m)
}

func resourcePreheatInstanceDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	_, _, respCode, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if respCode != 404 && err != nil { // We can't delete something that doesn't exist. Hence the 404-check
		return err
	}
	return nil
}
