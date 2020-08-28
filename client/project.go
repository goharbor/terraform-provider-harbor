package client

import (
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ProjectBody return a json body
func ProjectBody(d *schema.ResourceData) models.ProjectsBodyPost {
	body := models.ProjectsBodyPost{
		ProjectName: d.Get("name").(string),
	}
	body.Metadata.AutoScan = d.Get("vulnerability_scanning").(string)
	body.Metadata.Public = d.Get("public").(string)
	return body
}
