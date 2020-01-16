package provider

import (
	"encoding/json"
	"fmt"
	"strconv"

	"bitbucket.org/bestsellerit/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var pathProjects = "/api/projects"

type project struct {
	ProjectName string   `json:"project_name"`
	Metadata    metadata `json:"metadata"`
}

type metadata struct {
	AutoScan string `json:"auto_scan"`
	Public   string `json:"public"`
}

type projects struct {
	Name      string `json:"name"`
	ProjectID int    `json:"project_id"`
}

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"public": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  false,
			},
			"vulnerability_scanning": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  true,
			},
		},
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := project{
		ProjectName: d.Get("name").(string),
		Metadata: metadata{
			AutoScan: d.Get("vulnerability_scanning").(string),
			Public:   d.Get("public").(string),
		},
	}

	_, err := apiClient.SendRequest("POST", pathProjects, body, 0)
	if err != nil {
		return err
	}

	d.SetId(randomString(15))
	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := project{
		ProjectName: d.Get("name").(string),
	}

	resp, err := apiClient.SendRequest("GET", pathProjects+"?name="+body.ProjectName, nil, 0)
	if err != nil {
		return err
	}

	var jsonData []projects

	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to unmarchal: %s", err)
	}

	if len(jsonData) < 1 {
		return fmt.Errorf("[ERROR] JsonData is empty")
	}

	d.Set("project_id", jsonData[0].ProjectID)
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := project{
		ProjectName: d.Get("name").(string),
		Metadata: metadata{
			AutoScan: d.Get("vulnerability_scanning").(string),
			Public:   d.Get("public").(string),
		},
	}

	_, err := apiClient.SendRequest("PUT", pathProjects, body, 0)
	if err != nil {
		return err
	}

	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	id := d.Get("project_id").(int)

	apiClient.SendRequest("DELETE", pathProjects+"/"+strconv.Itoa(id), nil, 0)
	return nil
}
