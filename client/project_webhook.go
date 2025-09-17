package client

import (
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ProjectWebhookBody(d *schema.ResourceData) models.ProjectWebhook {
	eventTypes := d.Get("events_types").([]interface{})

	body := models.ProjectWebhook{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Enabled:     d.Get("enabled").(bool),
		EventTypes:  eventTypes,
	}
	targets := models.WebHookTargets{
		Type:           d.Get("notify_type").(string),
		AuthHeader:     d.Get("auth_header").(string),
		SkipCertVerify: d.Get("skip_cert_verify").(bool),
		Address:        d.Get("address").(string),
		PayloadFormat:  d.Get("payload_format").(string),
	}

	body.Targets = append(body.Targets, targets)

	return body
}
