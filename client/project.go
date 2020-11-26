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
		ProjectName: d.Get("name").(string),
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

	cveWhiteList := d.Get("cve_whitelist").([]interface{})
	log.Printf("[DEBUG] %v ", cveWhiteList)
	if len(cveWhiteList) > 0 {
		log.Printf("[DEBUG] %v ", expandCveWhiteList(cveWhiteList))
		body.CveWhitelist.Items = expandCveWhiteList(cveWhiteList)
		body.Metadata.ReuseSysCveWhitelist = "false"
	} else {
		body.Metadata.ReuseSysCveWhitelist = "true"
	}
	return body
}

func expandCveWhiteList(cveWhitelist []interface{}) models.CveWhitelistItems {
	whitelist := models.CveWhitelistItems{}

	for _, data := range cveWhitelist {
		item := models.CveWhitelistItem{
			CveID: data.(string),
		}
		whitelist = append(whitelist, item)
	}

	return whitelist
}
