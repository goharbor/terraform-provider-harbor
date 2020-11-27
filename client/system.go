package client

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetSystemBoby(d *schema.ResourceData, scheduleType string) models.SystemBody {
	var schedule string

	if scheduleType == "gc" {
		schedule = d.Get("schedule").(string)
	} else if scheduleType == "vuln" {
		schedule = d.Get("vulnerability_scan_policy").(string)
	}

	TypeStr, CronStr := GetSchedule(schedule)

	body := models.SystemBody{}
	body.Schedule.Type = TypeStr
	body.Schedule.Cron = CronStr
	if scheduleType == "gc" {
		body.Parameters.DeleteUntagged = d.Get("delete_untagged").(bool)
	}

	return body
}

// SetSchedule sets the schedule time to perform Vuln scanning and GC
func (client *Client) SetSchedule(d *schema.ResourceData, scheduleType string) (err error) {
	var path string

	if scheduleType == "gc" {
		path = models.PathGC
	} else if scheduleType == "vuln" {
		path = models.PathVuln
	}

	body := GetSystemBoby(d, scheduleType)

	resp, _, err := client.SendRequest("GET", path, nil, 200)
	if err != nil {
		return err
	}

	var jsonData models.SystemBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return err
	}

	time := jsonData.Schedule.Type
	requestType := "POST"
	httpStatusCode := 201

	if time != "" {
		log.Printf("Schedule found performing PUT request")
		requestType = "PUT"
		httpStatusCode = 200
	} else {
		log.Printf("No Schedule found performing POST request")
	}

	_, _, err = client.SendRequest(requestType, path, body, httpStatusCode)
	if err != nil {
		return err

	}
	return nil
}

// SetDefaultScanner set the default scanner within harbor
func (client *Client) SetDefaultScanner(scanner string) (err error) {
	resp, _, err := client.SendRequest("GET", models.PathScanners, nil, 0)

	body := models.ScannerBody{
		IsDefault: true,
	}

	var jsonData []models.ScannerBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return err
	}

	for _, v := range jsonData {

		if v.Name == strings.Title(scanner) {
			_, _, err = client.SendRequest("PATCH", models.PathScanners+"/"+v.UUID, body, 0)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
