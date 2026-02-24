package provider

import (
	"encoding/json"
	"fmt"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		ValidateRawResourceConfigFuncs: []schema.ValidateRawResourceConfigFunc{
			validation.PreferWriteOnlyAttribute(cty.GetAttrPath("password"), cty.GetAttrPath("password_wo")),
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ExactlyOneOf: []string{
					"password",
					"password_wo",
				},
				ConflictsWith: []string{
					"password_wo_version",
				},
			},
			"password_wo": {
				Type:      schema.TypeString,
				Optional:  true,
				WriteOnly: true,
				ExactlyOneOf: []string{
					"password",
					"password_wo",
				},
				RequiredWith: []string{
					"password_wo_version",
				},
			},
			"password_wo_version": {
				Type:     schema.TypeInt,
				Optional: true,
				RequiredWith: []string{
					"password_wo",
				},
				ConflictsWith: []string{
					"password",
				},
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
	passwordWriteOnly, err := getWriteOnlyString(d, "password_wo")
	if err != nil {
		return err
	}
	if passwordWriteOnly != "" {
		body.Password = passwordWriteOnly
		body.Newpassword = passwordWriteOnly
	}

	if body.Password == "" {
		return fmt.Errorf("one of password or password_wo must be configured")
	}

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
	oldPassword, _ := d.GetChange("password")
	oldPasswordWOVersion, _ := d.GetChange("password_wo_version")

	body := client.UserBody(d)
	passwordWriteOnly, err := getWriteOnlyString(d, "password_wo")
	if err != nil {
		return err
	}
	if passwordWriteOnly != "" {
		body.Password = passwordWriteOnly
		body.Newpassword = passwordWriteOnly
	}

	_, _, _, err = apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		return err
	}

	_, _, _, err = apiClient.SendRequest("PUT", d.Id()+"/sysadmin", body, 200)
	if err != nil {
		return err
	}

	if d.HasChange("password") || d.HasChange("password_wo_version") {
		if d.HasChange("password_wo_version") && passwordWriteOnly == "" && !d.HasChange("password") {
			_ = d.Set("password_wo_version", oldPasswordWOVersion)
			return fmt.Errorf("password_wo must be configured when password_wo_version changes")
		}
		_, _, _, err = apiClient.SendRequest("PUT", d.Id()+"/password", body, 200)
		if err != nil {
			if d.HasChange("password") {
				_ = d.Set("password", oldPassword)
			}
			if d.HasChange("password_wo_version") {
				_ = d.Set("password_wo_version", oldPasswordWOVersion)
			}
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
