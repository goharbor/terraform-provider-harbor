package provider

import (
	"encoding/json"
	"fmt"

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
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.UserBody(d)

	_, header, _, err := apiClient.SendRequest("POST", models.PathUsers, &body, 201)
	if err != nil {
		return err
	}

	id, err := client.GetID(header)
	if err != nil {
		return nil
	}

	d.SetId(id)
	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	resp, _, _, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	if err != nil {
		return err
	}
	var jsonData models.UserBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("Resource not found %s", d.Id())
	}

	d.Set("username", jsonData.Username)
	d.Set("full_name", jsonData.Realname)
	d.Set("email", jsonData.Email)
	d.Set("admin", jsonData.SysadminFlag)
	d.Set("comment", jsonData.Comment)

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.UserBody(d)
	_, _, _, err := apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		return err
	}

	_, _, _, err = apiClient.SendRequest("PUT", d.Id()+"/sysadmin", body, 200)
	if err != nil {
		return err
	}

	if d.HasChange("password") == true {
		_, _, _, err = apiClient.SendRequest("PUT", d.Id()+"/password", body, 200)
		if err != nil {
			return err
		}
	}

	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	_, _, _, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if err != nil {
		return err
	}
	return nil
}
