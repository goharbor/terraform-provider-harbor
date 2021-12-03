package provider

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRobotAccount() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"robot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"level": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"disable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"full_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"permissions": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access": {
							Type: schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:     schema.TypeString,
										Required: true,
									},
									"resource": {
										Type:     schema.TypeString,
										Required: true,
									},
									"effect": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "allow",
									},
								},
							},
							Required: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Required: true,
				ForceNew: true,
			},
		},
		Create: resourceRobotAccountCreate,
		Read:   resourceRobotAccountRead,
		Update: resourceRobotAccountUpdate,
		Delete: resourceRobotAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

	body := client.RobotBody(d)

	resp, headers, err := apiClient.SendRequest("POST", "/robots", body, 201)
	if err != nil {
		return err
	}

	var jsonData models.RobotBodyResponse
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return err
	}

	id, err := client.GetID(headers)
	if err != nil {
		return err
	}

	if d.Get("secret").(string) != "" {
		robotID := strconv.Itoa(jsonData.ID)
		secret := models.RobotSecret{
			Secret: d.Get("secret").(string),
		}
		_, _, err := apiClient.SendRequest("PATCH", "/robots/"+robotID, secret, 200)
		if err != nil {
			return err
		}
	} else {
		d.Set("secret", jsonData.Secret)
	}

	d.SetId(id)
	return resourceRobotAccountRead(d, m)
}

func resourceRobotAccountRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	robot, err := getRobot(d, apiClient)
	if err != nil {
		return err
	}

	d.Set("robot_id", strconv.Itoa(robot.ID))
	d.Set("full_name", robot.Name)

	return nil
}

func resourceRobotAccountUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.RobotBody(d)

	// if name not changed, use robot account name from api, otherwise it would always trigger a recreation,
	// since harbor does internally attach the robot account prefix to it´s names
	if false == d.HasChange("name") {
		robot, err := getRobot(d, apiClient)
		if err != nil {
			return err
		}
		body.Name = robot.Name
	}

	_, _, err := apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		return err
	}

	if d.HasChange("secret") {
		secret := models.RobotSecret{
			Secret: d.Get("secret").(string),
		}
		_, _, err := apiClient.SendRequest("PATCH", d.Id(), secret, 200)
		if err != nil {
			return err
		}
	}

	return resourceRobotAccountRead(d, m)
}

func resourceRobotAccountDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	_, _, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if err != nil {
		return err
	}
	return nil
}

func getRobot(d *schema.ResourceData, apiClient *client.Client) (models.RobotBody, error) {
	resp, _, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	if err != nil {
		return models.RobotBody{}, err
	}
	var jsonData models.RobotBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return models.RobotBody{}, fmt.Errorf("Resource not found %s", d.Id())
	}
	return jsonData, nil
}
