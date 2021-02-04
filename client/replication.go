package client

import (
	"log"

	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetReplicationBody(d *schema.ResourceData) models.ReplicationBody {

	action := d.Get("action").(string)
	schedule := d.Get("schedule").(string)

	body := models.ReplicationBody{
		Name:          d.Get("name").(string),
		Override:      d.Get("override").(bool),
		Enabled:       d.Get("enabled").(bool),
		DestNamespace: d.Get("dest_namespace").(string),
	}

	if action == "push" {
		body.DestRegistry.ID = d.Get("registry_id").(int)
	} else if action == "pull" {
		body.SrcRegistry.ID = d.Get("registry_id").(int)
	}

	switch schedule {
	case "manual":
		body.Trigger.Type = schedule

		break
	case "event_based":
		body.Trigger.Type = schedule
		break
	default:
		body.Trigger.Type = "schedule"
		body.Trigger.TriggerSettings.Cron = schedule
	}

	filters := d.Get("filters").(*schema.Set).List()

	for _, data := range filters {
		data := data.(map[string]interface{})
		filter := models.ReplicationFilters{}

		name := data["name"].(string)
		tag := data["tag"].(string)
		label := data["label"].(string)
		resource := data["resource"].(string)

		if name != "" {
			filter.Type = "name"
			filter.Value = name
		}
		if tag != "" {
			filter.Type = "tag"
			filter.Value = tag
		}
		if label != "" {
			filter.Type = "label"
			filter.Value = label
		}
		if resource != "" {
			filter.Type = "resource"
			filter.Value = resource
		}
		body.Filters = append(body.Filters, filter)

	}

	log.Println(body)
	return body
}
