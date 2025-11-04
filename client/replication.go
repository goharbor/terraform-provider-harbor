package client

import (
	"fmt"
	"log"
	"strings"

	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetReplicationBody(d *schema.ResourceData) models.ReplicationBody {

	action := d.Get("action").(string)
	schedule := d.Get("schedule").(string)

	body := models.ReplicationBody{
		Name:                    d.Get("name").(string),
		Description:             d.Get("description").(string),
		Override:                d.Get("override").(bool),
		Enabled:                 d.Get("enabled").(bool),
		Deletion:                d.Get("deletion").(bool),
		DestNamespace:           d.Get("dest_namespace").(string),
		DestNamespaceReplace:    d.Get("dest_namespace_replace").(int),
		CopyByChunk:             d.Get("copy_by_chunk").(bool),
		SingleActiveReplication: d.Get("single_active_replication").(bool),
		Speed:                   d.Get("speed").(int),
	}

	if action == "push" {
		body.DestRegistry.ID = d.Get("registry_id").(int)
	} else if action == "pull" {
		body.SrcRegistry.ID = d.Get("registry_id").(int)
	}

	switch schedule {
	case "manual", "event_based":
		body.Trigger.Type = schedule
		break
	default:
		body.Trigger.Type = "scheduled"
		body.Trigger.TriggerSettings.Cron = schedule
	}

	filters := d.Get("filters").(*schema.Set).List()

	for _, data := range filters {
		data := data.(map[string]interface{})
		filter := models.ReplicationFilters{}

		name := data["name"].(string)
		tag := data["tag"].(string)
		label := data["labels"].([]interface{})
		decoration := data["decoration"].(string)
		resource := data["resource"].(string)

		if name != "" {
			filter.Type = "name"
			filter.Value = name
		}
		if tag != "" {
			filter.Type = "tag"
			filter.Value = strings.ReplaceAll(tag, " ", "")
			filter.Decoration = decoration
		}
		if len(label) > 0 {
			filter.Type = "label"
			filter.Value = make([]string, 0)
			for _, v := range label {
				filter.Value = append(filter.Value.([]string), fmt.Sprintf("%v", v))
			}
			filter.Decoration = decoration
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

func GetExecutionBody(d *schema.ResourceData) models.ExecutionBody {

	body := models.ExecutionBody{
		PolicyID: d.Get("policy_id").(int),
	}

	log.Println(body)
	return body
}
