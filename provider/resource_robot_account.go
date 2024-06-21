package provider

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
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
	if !strings.Contains(id, path) {
		id = path + id
	}
	return id

}

func resourceRobotAccountCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.RobotBody(d)

	resp, headers, _, err := apiClient.SendRequest("POST", "/robots", body, 201)
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
		_, _, _, err := apiClient.SendRequest("PATCH", "/robots/"+robotID, secret, 200)
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
		d.SetId("")
		return nil
	}

	resp, _, respCode, err := apiClient.SendRequest("GET", models.PathConfig, nil, 200)
	if respCode == 404 && err != nil {
		d.SetId("")
		return fmt.Errorf("error getting system configuration %s", err)
	}
	var systemConfig models.ConfigBodyResponse
	err = json.Unmarshal([]byte(resp), &systemConfig)
	if err != nil {
		return fmt.Errorf("error getting system configuration %s", err)
	}

	shortName := strings.TrimPrefix(robot.Name, systemConfig.RobotNamePrefix.Value)
	if robot.Level == "project" {
		shortName = strings.Split(shortName, "+")[1]
	}

	d.Set("name", shortName)
	d.Set("robot_id", strconv.Itoa(robot.ID))
	d.Set("full_name", robot.Name)
	d.Set("description", robot.Description)
	d.Set("level", robot.Level)
	d.Set("duration", robot.Duration)
	d.Set("disable", robot.Disable)

	// Set the permissions of the robot account in the Terraform state
	permissions := make([]map[string]interface{}, len(robot.Permissions))
	for i, permission := range robot.Permissions {
		permissionMap := make(map[string]interface{})
		permissionMap["kind"] = permission.Kind
		permissionMap["namespace"] = permission.Namespace
		access := make([]map[string]interface{}, len(permission.Access))
		for i, v := range permission.Access {
			accessMap := make(map[string]interface{})
			accessMap["action"] = v.Action
			accessMap["resource"] = v.Resource
			accessMap["effect"] = v.Effect
			access[i] = accessMap
		}
		permissionMap["access"] = access
		permissions[i] = permissionMap
	}
	d.Set("permissions", permissions)
	return nil
}

func resourceRobotAccountUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.RobotBody(d)

	// if name not changed, use robot account name from api, otherwise it would always trigger a recreation,
	// since harbor does internally attach the robot account prefix to it´s names
	if !d.HasChange("name") {
		robot, err := getRobot(d, apiClient)
		if err != nil {
			return err
		}
		body.Name = robot.Name
	}

	_, _, _, err := apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		return err
	}

	if d.HasChange("secret") {
		secret := models.RobotSecret{
			Secret: d.Get("secret").(string),
		}
		_, _, _, err := apiClient.SendRequest("PATCH", d.Id(), secret, 200)
		if err != nil {
			return err
		}
	}

	return resourceRobotAccountRead(d, m)
}

func resourceRobotAccountDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	_, _, _, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if err != nil {
		return err
	}
	return nil
}

func getRobot(d *schema.ResourceData, apiClient *client.Client) (models.RobotBody, error) {
	resp, _, _, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	if err != nil {
		return models.RobotBody{}, err
	}
	var jsonData models.RobotBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return models.RobotBody{}, fmt.Errorf("resource not found %s", d.Id())
	}
	return jsonData, nil
}
