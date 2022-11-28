package client

import (
	"log"

	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SystemCVEAllowListBody(d *schema.ResourceData) models.SystemCveAllowListBodyPost {
	body := models.SystemCveAllowListBodyPost{
		ProjectID: d.Get("project_id").(string),
	}

	expires_at, expires_at_true := d.GetOk("expires_at")

	if expires_at_true {
		body.ExpiresAt = expires_at.(int)
	}

	cveAllowList := d.Get("cve_allowlist").([]interface{})
	log.Printf("[DEBUG] %v ", cveAllowList)
	if len(cveAllowList) > 0 {
		log.Printf("[DEBUG] %v ", expandSystemCveAllowList(cveAllowList))
		body.Items = expandSystemCveAllowList(cveAllowList)
	}

	return body
}

func expandSystemCveAllowList(cveAllowlist []interface{}) models.SystemCveAllowlistItems {
	allowlist := models.SystemCveAllowlistItems{}

	for _, data := range cveAllowlist {
		item := models.SystemCveAllowlistItem{
			CveID: data.(string),
		}
		allowlist = append(allowlist, item)
	}

	return allowlist
}
