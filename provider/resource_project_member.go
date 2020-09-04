package provider

import (
	"encoding/json"
	"fmt"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMembers() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{

			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
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
					if v != "projectadmin" && v != "developer" && v != "guest" && v != "master" {
						errs = append(errs, fmt.Errorf("%q must be either projectadmin, developer, guest or master, got: %s", key, v))
					}
					return
				},
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "ldap" && v != "internal" && v != "oidc" {
						errs = append(errs, fmt.Errorf("%q must be either ldap, internal or oidc, got: %s", key, v))
					}
					return
				},
			},
		},
		Create:             resourceMembersCreate,
		Read:               resourceMembersRead,
		Update:             resourceMembersUpdate,
		Delete:             resourceMembersDelete,
		DeprecationMessage: "The resource project_member has been renamed to project_member_group. This resource is deprecated and will be removed in the next major version",
	}
}

func resourceMembersCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	projectid := checkProjectid(d.Get("project_id").(string))
	path := projectid + "/members"

	body := client.ProjectMembersGroupBody(d)

	_, headers, err := apiClient.SendRequest("POST", path, body, 201)
	if err != nil {
		return err
	}

	id, err := client.GetID(headers)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceMembersRead(d, m)
}

func resourceMembersRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	if err != nil {
		fmt.Println(err)
	}

	var jsonData models.ProjectMembersBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return err
	}

	d.Set("role", client.RoleTypeNumber(jsonData.RoleID))
	return nil
}

func resourceMembersUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.ProjectMembersGroupBody(d)
	_, _, err := apiClient.SendRequest("GET", d.Id(), body, 200)
	if err != nil {
		fmt.Println(err)
	}

	return resourceMembersRead(d, m)
}

func resourceMembersDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	_, _, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
