package provider

import (
	"encoding/json"
	"log"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var pathVuln = "/system/scanAll/schedule"
var TypeStr string
var CronStr string

type schedule struct {
	Schedule cron `json:"schedule`
}

type cron struct {
	Type string `json:"type"`
	Cron string `json:"cron`
}

type Schedule2 struct {
	Type string `json:"type"`
	Cron string `json:"cron"`
}
type Info struct {
	Schedule Schedule2 `json:schedule`
}

func resourceTasks() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Use harbor_interrogation_services resource instead. harbor_tasks Will be removed in the next major version",
		Schema: map[string]*schema.Schema{
			"vulnerability_scan_policy": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourceTasksCreate,
		Read:   resourceTasksRead,
		Update: resourceTasksUpdate,
		Delete: resourceTasksDelete,
	}
}

func resourceTasksCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	vulnSchedule := d.Get("vulnerability_scan_policy").(string)
	getSchedule(vulnSchedule)

	body := schedule{
		Schedule: cron{
			Type: TypeStr,
			Cron: CronStr,
		},
	}

	resp, _, _, err := apiClient.SendRequest("GET", pathVuln, nil, 0)
	if err != nil {
		return err
	}
	var jsonData Info

	json.Unmarshal([]byte(resp), &jsonData)

	time := jsonData.Schedule.Type
	requestType := "POST"
	if time != "" {
		log.Printf("Schedule found performing PUT request")
		requestType = "PUT"
	} else {
		log.Printf("No schedule found performing POST request")
	}
	_, _, _, err = apiClient.SendRequest(requestType, pathVuln, body, 0)
	if err != nil {
		return err

	}

	d.SetId(randomString(15))
	return resourceTasksRead(d, m)
}

func resourceTasksRead(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceTasksUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	vulnSchedule := d.Get("vulnerability_scan_policy").(string)
	getSchedule(vulnSchedule)

	body := schedule{
		Schedule: cron{
			Type: TypeStr,
			Cron: CronStr,
		},
	}

	_, _, _, err := apiClient.SendRequest("PUT", pathVuln, body, 0)
	if err != nil {
		return err
	}
	return resourceTasksRead(d, m)
}

func resourceTasksDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := schedule{
		Schedule: cron{Cron: ""},
	}

	_, _, _, err := apiClient.SendRequest("PUT", pathVuln, body, 0)
	if err != nil {
		return err
	}
	return nil
}

func getSchedule(schedule string) {
	switch schedule {
	case "hourly":
		TypeStr = "Hourly"
		CronStr = "0 0 * * * *"
	case "daily":
		TypeStr = "Daily"
		CronStr = "0 0 0 * * *"
	case "weekly":
		TypeStr = "Weekly"
		CronStr = "0 0 0 * * 0"
	}
}
