package client

import (
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func LabelsBody(d *schema.ResourceData) models.Labels {
	body := models.Labels{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Color:       d.Get("color").(string),
	}

	if _, ok := d.GetOk("project_id"); ok {

		body.ID = d.Get("project_id").(int)
		body.Scope = "p"
	} else {
		body.Scope = "g"
	}
	return body
}

// ) models.ProjectsBodyPost {
// 	body := models.ProjectsBodyPost{
// 		ProjectName: d.Get("name").(string),
// 	}
// 	body.Metadata.AutoScan = d.Get("vulnerability_scanning").(string)
// 	body.Metadata.Public = d.Get("public").(string)
// 	return body
