package provider

import (
	"encoding/json"
	"fmt"
	"strconv"

	"bitbucket.org/bestsellerit/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var pathRobot string = "/api/projects"

type robot struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Access      []access `json:"access"`
}

type access struct {
	Action   string `json:"action"`
	Resource string `json:"resource"`
}

type robotAccount struct {
	Token   string `json:"token,omitempty"`
	RobotID int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
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
	}
}

func resourceRobotAccountCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	projectid := d.Get("project_id").(string)

	url := pathRobot + "/" + projectid + "/robots"
	resource := "/project/" + projectid + "/repository"

	body := robot{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Access: []access{
			access{
				Action:   d.Get("action").(string),
				Resource: resource,
			},
		},
	}

	resp, err := apiClient.SendRequest("POST", url, body, 201)
	if err != nil {
		return err
	}

	var jsonData robotAccount

	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to unmarshal: %s", err)
	}

	d.Set("token", jsonData.Token)
	return resourceRobotAccountRead(d, m)
}

func resourceRobotAccountRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	projectid := d.Get("project_id").(string)
	name := d.Get("name").(string)
	url := pathRobot + "/" + projectid + "/robots"

	resp, err := apiClient.SendRequest("GET", url, nil, 200)
	if err != nil {
		return err
	}

	var jsonData []robotAccount

	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to unmarchal: %s", err)
	}

	if len(jsonData) < 1 {
		return fmt.Errorf("[ERROR] JsonData is empty")
	}

	for _, v := range jsonData {
		if v.Name == "robot$"+name {
			d.SetId(strconv.Itoa(v.RobotID))
			d.Set("robot_id", strconv.Itoa(v.RobotID))
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
	projectid := d.Get("project_id").(string)
	robotid := d.Get("robot_id").(string)
	url := pathRobot + "/" + projectid + "/robots/" + robotid
	apiClient.SendRequest("DELETE", url, nil, 0)

	return nil
}
