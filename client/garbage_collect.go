package client

import (
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func GCBodyPost(d *schema.ResourceData) models.GCBodyPost {
	body := models.GCBodyPost{}
	body.Schedule.Type = d.Get("type").(string)
	body.Parameters.DeleteUntagged = d.Get("delete_untagged").(bool)
	
	if body.Schedule.Type == "Hourly" {
		body.Schedule.Cron = "0 0 * * * *"
	}	else if body.Schedule.Type == "Daily" {
		body.Schedule.Cron = "0 0 0 * * *"
	}	else if body.Schedule.Type == "Weekly" {
		body.Schedule.Cron = "0 0 * * * 0"
	} else if body.Schedule.Type == "Custom" {
		body.Schedule.Cron = d.Get("cron").(string)
	} else if body.Schedule.Type == "None" {
		body.Schedule.Cron = ""
	}

	return body
}

func GCBodyDelete() models.GCBodyPost {
	body := models.GCBodyPost{}
	body.Schedule.Type = "None"
	body.Schedule.Cron = ""
	body.Parameters.DeleteUntagged = false
	return body
}