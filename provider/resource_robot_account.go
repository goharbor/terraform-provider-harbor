package provider

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"bitbucket.org/bestsellerit/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// var pathRobot string = "/projects"

type robotRequest struct {
	Access      []access `json:"access,omitempty"`
	Name        string   `json:"name,omitempty"`
	ExpiresAt   int      `json:"expires_at,omitempty"`
	Description string   `json:"description,omitempty"`
}
type access struct {
	Action   string `json:"action,omitempty"`
	Resource string `json:"resource,omitempty"`
}

type robotRepones struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Token       string `json:"token"`
	Description string `json:"description"`
	ProjectID   int    `json:"project_id"`
	ExpiresAt   int    `json:"expires_at"`
	Disabled    bool   `json:"disabled"`
}

func resourceRobotAccount() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "pull",
				ForceNew: true,
			},
			"token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"robot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Create: resourceRobotAccountCreate,
		Read:   resourceRobotAccountRead,
		Update: resourceRobotAccountUpdate,
		Delete: resourceRobotAccountDelete,
		Importer: &schema.ResourceImporter{
			// State: schema.ImportStatePassthrough,
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

				apiclient := m.(*client.Client)

				resp, err := apiclient.SendRequest("GET", d.Id(), nil, 200)
				if err != nil {
					fmt.Println(err)
				}
				var jsonData robotRepones
				json.Unmarshal([]byte(resp), &jsonData)

				d.Set("name", strings.Replace(jsonData.Name, "robot$", "", -1))
				d.Set("description", jsonData.Description)
				// d.Set("prjoect", "projects/"+jsonData.ProjectID)

				// d.SetId(d.Id())

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func checkProjectid(id string) (projecid string) {
	path := "/projects/"
	if strings.Contains(id, path) == false {
		id = path + id
	}
	return id

}

func resourceRobotAccountCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	projectid := checkProjectid(d.Get("project_id").(string))

	url := projectid + "/robots"
	// action := "Action:" + d.Get("action").(string)
	resource := strings.Replace(projectid, "s", "", +1) + "/repository"

	body := robotRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Access: []access{
			{
				Action:   d.Get("action").(string),
				Resource: resource,
			},
		},
	}

	resp, err := apiClient.SendRequest("POST", url, body, 201)
	if err != nil {
		return err
	}

	var jsonData robotRepones

	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to unmarshal: %s", err)
	}

	d.Set("token", jsonData.Token)
	return resourceRobotAccountRead(d, m)
}

func resourceRobotAccountRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	projectid := checkProjectid(d.Get("project_id").(string))
	name := d.Get("name").(string)

	url := projectid + "/robots"

	var jsonData []robotRepones

	resp, err := apiClient.SendRequest("GET", url, nil, 200)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to unmarchal: %s", err)
	}

	if len(jsonData) < 1 {
		return fmt.Errorf("[ERROR] JsonData is empty")
	}

	for _, v := range jsonData {
		if v.Name == "robot$"+name {
			d.SetId(url + "/" + strconv.Itoa(v.ID))
			d.Set("robot_id", strconv.Itoa(v.ID))
			d.Set("name", strings.Replace(name, "robot$", "", -1))
			d.Set("description", v.Description)
		}
	}
	return nil
}

func resourceRobotAccountUpdate(d *schema.ResourceData, m interface{}) error {
	// apiClient := m.(*client.Client)

	return resourceRobotAccountRead(d, m)
}

func resourceRobotAccountDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	apiClient.SendRequest("DELETE", d.Id(), nil, 200)

	return nil
}
