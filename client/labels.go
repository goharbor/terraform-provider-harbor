package client

import (
	"strconv"
	"strings"

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
		id := d.Get("project_id").(string)

		body.ID, _ = strconv.Atoi(strings.ReplaceAll(id, "/projects/", ""))
		body.Scope = "p"
	} else {

		body.Scope = "g"
	}
	return body
}
