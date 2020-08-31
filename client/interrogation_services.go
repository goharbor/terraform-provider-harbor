package client

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func GetVulnBoby(d *schema.ResourceData) models.VulnBody {
	vulnSchedule := d.Get("vulnerability_scan_policy").(string)
	TypeStr, CronStr := GetSchedule(vulnSchedule)

	body := models.VulnBody{}
	body.Schedule.Type = TypeStr
	body.Schedule.Cron = CronStr
	return body
}

// SetScannerPolicy sets the schedule time to perform Vuln scannin
func (client *Client) SetScannerPolicy(d *schema.ResourceData) (err error) {

	body := GetVulnBoby(d)

	resp, _, err := client.SendRequest("GET", models.PathVuln, nil, 0)
	if err != nil {
		return err
	}

	var jsonData models.VulnBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return err
	}

	time := jsonData.Schedule.Type
	requestType := "POST"
	if time != "" {
		log.Printf("Shedule found performing PUT request")
		requestType = "PUT"
	} else {
		log.Printf("No shedule found performing POST request")
	}

	_, _, err = client.SendRequest(requestType, models.PathVuln, body, 200)
	if err != nil {
		return err

	}
	return nil
}

// SetDefaultScanner set the default scanner within harbor
func (client *Client) SetDefaultScanner(scanner string) (err error) {
	resp, _, err := client.SendRequest("GET", models.PathScanners, nil, 200)

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
			_, _, err = client.SendRequest("PATCH", models.PathScanners+"/"+v.UUID, body, 200)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
