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

func dataRobotAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataRobotAccountsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"level": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"robot_accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataRobotAccountsRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	page := 1
	name := d.Get("name").(string)
	level := d.Get("level").(string)
	projectId := d.Get("project_id").(int)

	robotsQueryPath := []string{}
	if level != "" {
		robotsQueryPath = append(robotsQueryPath, "Level="+level)
	}
	if projectId != 0 {
		robotsQueryPath = append(robotsQueryPath, "ProjectID="+strconv.Itoa(projectId))
	}

	robotAccountsData := make([]map[string]interface{}, 0)
	for {
		robotsPath := models.PathRobots + "?page=" + strconv.Itoa(page)
		if len(robotsQueryPath) > 0 {
			robotsPath += "&q=" + strings.Join(robotsQueryPath, ",")
		}

		resp, _, _, err := apiClient.SendRequest("GET", robotsPath, nil, 200)
		if err != nil {
			return err
		}

		var jsonData []models.RobotBody
		err = json.Unmarshal([]byte(resp), &jsonData)
		if err != nil {
			return fmt.Errorf("unable to retrieve Harbor robot accounts data: %s", err)
		}

		// If there is no data on the current page, we have reached the last page
		if len(jsonData) == 0 {
			break
		}

		for _, v := range jsonData {
			if name == "" || v.Name == "robot$"+name {
				id := models.PathRobots + "/" + strconv.Itoa(v.ID)

				robotAccountData := map[string]interface{}{
					"id":          id,
					"name":        v.Name,
					"description": v.Description,
					"level":       v.Level,
					"duration":    v.Duration,
					"disable":     v.Disable,
				}

				robotAccountsData = append(robotAccountsData, robotAccountData)
			}
		}

		page++
	}
	d.SetId("harbor-robot-accounts")
	d.Set("robot_accounts", robotAccountsData)

	return nil
}
