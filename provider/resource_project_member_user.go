package provider

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMembersUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"member_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validRoles := []string{"projectadmin", "developer", "guest", "maintainer", "limitedguest"}
					for _, r := range validRoles {
						if v == r {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validRoles, v))
					return
				},
			},
		},
		Create: resourceMembersUserCreate,
		Read:   resourceMembersUserRead,
		Update: resourceMembersUserUpdate,
		Delete: resourceMembersUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceMembersUserCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	projectid := checkProjectid(d.Get("project_id").(string))
	path := projectid + "/members"

	body := client.ProjectMembersUserBody(d)

	_, headers, _, err := apiClient.SendRequest("POST", path, body, 201)
	if err != nil {
		return err
	}

	id, err := client.GetID(headers)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceMembersUserRead(d, m)
}

func resourceMembersUserRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, _, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	if err != nil {
		d.SetId("")
		return nil
	}

	var jsonData models.ProjectMembersBodyResponses
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("Resource not found %s", d.Id())
	}

	d.Set("role", client.RoleTypeNumber(jsonData.RoleID))
	d.Set("project_id", checkProjectid(strconv.Itoa(jsonData.ProjectID)))
	d.Set("user_name", jsonData.EntityName)
	return nil
}

func resourceMembersUserUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.ProjectMembersUserBody(d)
	_, _, _, err := apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		fmt.Println(err)
	}

	return resourceMembersUserRead(d, m)
}

func resourceMembersUserDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	_, _, _, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
