package provider

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMembersGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"group_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"group_id", "group_name"},
			},
			"ldap_group_dn": {
				Type:     schema.TypeString,
				Optional: true,
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
					if v != "projectadmin" && v != "developer" && v != "guest" && v != "master" && v != "limitedguest" {
						errs = append(errs, fmt.Errorf("%q must be either projectadmin, developer, guest, limitedguest or master, got: %s", key, v))
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
		Create: resourceMembersGroupCreate,
		Read:   resourceMembersGroupRead,
		Update: resourceMembersGroupUpdate,
		Delete: resourceMembersGroupDelete,
	}
}

func resourceMembersGroupCreate(d *schema.ResourceData, m interface{}) error {
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
	return resourceMembersGroupRead(d, m)
}

func resourceMembersGroupRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
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
	d.Set("project_id", strconv.Itoa(jsonData.ProjectID))
	d.Set("group_name", jsonData.EntityName)
	return nil
}

func resourceMembersGroupUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.ProjectMembersGroupBody(d)
	_, _, err := apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		fmt.Println(err)
	}

	return resourceMembersGroupRead(d, m)
}

func resourceMembersGroupDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	_, _, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
