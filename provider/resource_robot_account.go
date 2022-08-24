package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
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
		CreateContext: resourceRobotAccountCreate,
		ReadContext:   resourceRobotAccountRead,
		UpdateContext: resourceRobotAccountUpdate,
		DeleteContext: resourceRobotAccountDelete,
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

func resourceRobotAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	body := client.RobotBody(d)

	resp, headers, _, err := apiClient.SendRequest(ctx, "POST", "/robots", body, 201)
	if err != nil {
		return diag.FromErr(err)
	}

	var jsonData models.RobotBodyResponse
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := client.GetID(headers)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("secret").(string) != "" {
		robotID := strconv.Itoa(jsonData.ID)
		secret := models.RobotSecret{
			Secret: d.Get("secret").(string),
		}
		_, _, _, err := apiClient.SendRequest(ctx, "PATCH", "/robots/"+robotID, secret, 200)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		d.Set("secret", jsonData.Secret)
	}

	d.SetId(id)
	return resourceRobotAccountRead(ctx, d, m)
}

func resourceRobotAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	resp, _, respCode, err := apiClient.SendRequest(ctx, "GET", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Robot account %q was not found - removing from state!", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("making Read request on robot account %s : %+v", d.Id(), err)
	}

	var robot models.RobotBody
	err = json.Unmarshal([]byte(resp), &robot)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("robot_id", strconv.Itoa(robot.ID))
	d.Set("full_name", robot.Name)

	return nil
}

func resourceRobotAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	body := client.RobotBody(d)

	// if name not changed, use robot account name from api, otherwise it would always trigger a recreation,
	// since harbor does internally attach the robot account prefix to it´s names
	if !d.HasChange("name") {
		robot, err := getRobot(ctx, d, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		body.Name = robot.Name
	}

	_, _, _, err := apiClient.SendRequest(ctx, "PUT", d.Id(), body, 200)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("secret") {
		secret := models.RobotSecret{
			Secret: d.Get("secret").(string),
		}
		_, _, _, err := apiClient.SendRequest(ctx, "PATCH", d.Id(), secret, 200)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceRobotAccountRead(ctx, d, m)
}

func resourceRobotAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	_, _, respCode, err := apiClient.SendRequest(ctx, "DELETE", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Robot account %q was not found - already deleted!", d.Id())
			return nil
		}
		return diag.Errorf("making delete request on robot account %s : %+v", d.Id(), err)
	}
	return nil
}

func getRobot(ctx context.Context, d *schema.ResourceData, apiClient *client.Client) (models.RobotBody, error) {
	resp, _, _, err := apiClient.SendRequest(ctx, "GET", d.Id(), nil, 200)
	if err != nil {
		return models.RobotBody{}, fmt.Errorf("making read request on robot account %s : %+v", d.Id(), err)
	}

	if err != nil {
		return models.RobotBody{}, err
	}

	var jsonData models.RobotBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return models.RobotBody{}, err
	}

	return jsonData, nil
}
