package provider

import (
	"encoding/json"
	"log"

	"bitbucket.org/bestsellerit/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var pathVuln = "/api/system/scanAll/schedule"
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

	resp, _ := apiClient.SendRequest("GET", pathVuln, nil, 0)
	var jsonData Info

	json.Unmarshal([]byte(resp), &jsonData)

	time := jsonData.Schedule.Type
	if time != "" {
		log.Printf("Shedule found performing PUT request")
		apiClient.SendRequest("PUT", pathVuln, body, 0)
	} else {
		log.Printf("No shedule found performing POST request")
		apiClient.SendRequest("POST", pathVuln, body, 0)
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

	apiClient.SendRequest("PUT", pathVuln, body, 0)

	return resourceTasksRead(d, m)
}

func resourceTasksDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := schedule{
		Schedule: cron{Cron: ""},
	}

	apiClient.SendRequest("PUT", pathVuln, body, 0)
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
