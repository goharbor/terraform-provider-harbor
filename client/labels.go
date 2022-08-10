package client

import (
	"log"
	"strconv"
	"strings"

	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func LabelsBody(d *schema.ResourceData) models.Labels {
	body := models.Labels{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Color:       d.Get("color").(string),
	}

	project := d.Get("project_id").(string)
	if project != "" {
		id, err := strconv.Atoi(strings.ReplaceAll(project, "/projects/", ""))
		if err != nil {
			log.Println(err)
		}

		body.ProjectID = id
		body.Scope = "p"
	} else {

		body.Scope = "g"
	}
	return body
}
