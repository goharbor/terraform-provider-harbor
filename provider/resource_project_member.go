package provider

import (
	"encoding/json"
	"fmt"
	"strconv"

	"bitbucket.org/bestsellerit/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// var pathProjects = "/api/projects"

type members struct {
	RoleID      int   `json:"role_id"`
	GroupMember group `json:"member_group"`
}

type group struct {
	GroupType int    `json:"group_type"`
	GroupName string `json:"group_name"`
}

type entity struct {
	ID     int `json:"id"`
	RoleID int `json:"role_id"`
}

func resourceMembers() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeInt,
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
		Create: resourceMembersCreate,
		Read:   resourceMembersRead,
		Update: resourceMembersUpdate,
		Delete: resourceMembersDelete,
	}
}

func groupType(group string) (x int) {
	switch group {
	case "ldap":
		x = 1
	case "internal":
		x = 2
	case "oidc":
		x = 3
	}
	return x
}

func roleType(role string) (x int) {
	switch role {
	case "projectadmin":
		x = 1
	case "developer":
		x = 2
	case "guest":
		x = 3
	case "master":
		x = 4
	}
	return x
}

func resourceMembersCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	path := pathProjects + "/" + strconv.Itoa(d.Get("project_id").(int)) + "/members"

	body := members{
		RoleID: roleType(d.Get("role").(string)),
		GroupMember: group{
			GroupType: groupType(d.Get("type").(string)),
			GroupName: d.Get("name").(string),
		},
	}

	_, err := apiClient.SendRequest("POST", path, body, 201)
	if err != nil {
		return err
	}

	d.SetId(randomString(15))
	return resourceMembersRead(d, m)
}

func resourceMembersRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	path := pathProjects + "/" + strconv.Itoa(d.Get("project_id").(int)) + "/members?entityname=" + d.Get("name").(string)

	resp, err := apiClient.SendRequest("GET", path, nil, 200)
	if err != nil {
		fmt.Println(err)
	}

	var entityData []entity
	json.Unmarshal([]byte(resp), &entityData)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to unmarshal: %s", err)
	}
	fmt.Printf("%v", entityData)

	d.Set("member_id", entityData[0].ID)
	// d.Set("")
	return nil
}

func resourceMembersUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	path := pathProjects + "/" + strconv.Itoa(d.Get("project_id").(int)) + "/members/" + strconv.Itoa(d.Get("member_id").(int))

	_, err := apiClient.SendRequest("GET", path, nil, 200)
	if err != nil {
		fmt.Println(err)
	}

	return resourceMembersRead(d, m)
}

func resourceMembersDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	path := pathProjects + "/" + strconv.Itoa(d.Get("project_id").(int)) + "/members/" + strconv.Itoa(d.Get("member_id").(int))

	_, err := apiClient.SendRequest("DELETE", path, nil, 200)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
