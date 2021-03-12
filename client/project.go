package client

import (
	"log"
	"strconv"

	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProjectBody return a json body
func ProjectBody(d *schema.ResourceData) models.ProjectsBodyPost {
	body := models.ProjectsBodyPost{
		ProjectName:  d.Get("name").(string),
		RegistryID:   d.Get("registry_id").(int),
		StorageLimit: d.Get("storage_quota").(int) * 1073741824,
	}
	body.Metadata.AutoScan = strconv.FormatBool(d.Get("vulnerability_scanning").(bool))
	body.Metadata.Public = d.Get("public").(string)

	security := d.Get("deployment_security").(string)
	if security != "" {
		body.Metadata.Severity = security
		body.Metadata.PreventVul = "true"
	} else {
		body.Metadata.Severity = "low"
		body.Metadata.PreventVul = "false"
	}

	body.Metadata.EnableContentTrust = strconv.FormatBool(d.Get("enable_content_trust").(bool))

	cveAllowList := d.Get("cve_allowlist").([]interface{})
	log.Printf("[DEBUG] %v ", cveAllowList)
	if len(cveAllowList) > 0 {
		log.Printf("[DEBUG] %v ", expandCveAllowList(cveAllowList))
		body.CveAllowlist.Items = expandCveAllowList(cveAllowList)
		body.Metadata.ReuseSysCveAllowlist = "false"
	} else {
		body.Metadata.ReuseSysCveAllowlist = "true"
	}

	return body
}

func expandCveAllowList(cveAllowlist []interface{}) models.CveAllowlistItems {
	allowlist := models.CveAllowlistItems{}

	for _, data := range cveAllowlist {
		item := models.CveAllowlistItem{
			CveID: data.(string),
		}
		allowlist = append(allowlist, item)
	}

	return allowlist
}
