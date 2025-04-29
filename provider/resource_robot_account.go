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
										//Default:  "allow",
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

	resp, headers, _, err := apiClient.SendRequest("POST", models.PathRobots, body, 201)
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
		_, _, _, err := apiClient.SendRequest("PATCH", models.PathRobots+"/"+robotID, secret, 200)
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

	var shortName string

	if robot.Level == "project" {
		// if it's a project level robot account, we just need to split on the "+" char as it's not an allowed char in the project name and always the separator for robot project name
		// eg : robot name = "robot$project123+robot123"
		shortName = strings.Split(robot.Name, robot.Permissions[0].Namespace+"+")[1]
	} else {
		// if it's a system level robot account, we check if the robot_prefix is set
		if m.(*client.Client).GetRobotPrefix() != "" {
			// if robot_prefix is set, we use it to get the short name
			shortName = strings.TrimPrefix(robot.Name, m.(*client.Client).GetRobotPrefix())
		} else {
			// if robot_prefix is not set, we need to get the system configuration to get the prefix
			resp, _, respCode, err := apiClient.SendRequest("GET", models.PathConfig, nil, 200)
			if respCode == 404 && err != nil {
				d.SetId("")
				return fmt.Errorf("error getting system configuration (probably missing admin rights) %s, you can use robot_prefix to force the prefix", err)
			}
			var systemConfig models.ConfigBodyResponse
			err = json.Unmarshal([]byte(resp), &systemConfig)
			if err != nil {
				return fmt.Errorf("error getting system configuration (probably missing admin rights) %s, you can use robot_prefix to force the prefix", err)
			}
			shortName = strings.TrimPrefix(robot.Name, systemConfig.RobotNamePrefix.Value)
		}
	}
	d.Set("name", shortName)
	d.Set("robot_id", strconv.Itoa(robot.ID))
	d.Set("full_name", robot.Name)
	d.Set("description", robot.Description)
	d.Set("level", robot.Level)
	d.Set("duration", robot.Duration)
	d.Set("disable", robot.Disable)

	// Set the permissions of the robot account in the Terraform state
	tfPermissions := d.Get("permissions").(*schema.Set).List()
	permissions := make([]map[string]interface{}, len(robot.Permissions))
	for i, permission := range robot.Permissions {
		permissionMap := make(map[string]interface{})
		permissionMap["kind"] = permission.Kind
		permissionMap["namespace"] = permission.Namespace

		// Find matching tfvars permission by kind and namespace.
		tfPerm, found := findTfPermission(tfPermissions, permission.Kind, permission.Namespace)
		var tfAccessSet []interface{}
		if found {
			tfAccessSet = tfPerm["access"].(*schema.Set).List()
		} else {
			tfAccessSet = []interface{}{}
		}

		accessList := make([]map[string]interface{}, len(permission.Access))
		for j, v := range permission.Access {
			accessMap := make(map[string]interface{})
			accessMap["action"] = v.Action
			accessMap["resource"] = v.Resource

			// Search for a matching access block in tfAccessSet.
			if tfAccess, ok := findTfAccess(tfAccessSet, v.Action, v.Resource); ok {
				// If explicitly set to "allow" in tfvars, enforce that.
				if eff, exists := tfAccess["effect"]; exists && eff != "" {
					accessMap["effect"] = eff
				}
			} else {
				// If no matching tfvars block is found or API default is not "allow", use the API value if non-empty.
				if v.Effect != "" && v.Effect != "allow" {
					accessMap["effect"] = v.Effect
				}
			}
			accessList[j] = accessMap
		}
		permissionMap["access"] = accessList
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

func findTfPermission(tfPermissions []interface{}, kind, namespace string) (map[string]interface{}, bool) {
	for _, item := range tfPermissions {
		tfPerm := item.(map[string]interface{})
		if tfPerm["kind"] == kind && tfPerm["namespace"] == namespace {
			return tfPerm, true
		}
	}
	return nil, false
}

func findTfAccess(tfAccesses []interface{}, action, resource string) (map[string]interface{}, bool) {
	for _, item := range tfAccesses {
		access := item.(map[string]interface{})
		if access["action"] == action && access["resource"] == resource {
			return access, true
		}
	}
	return nil, false
}
