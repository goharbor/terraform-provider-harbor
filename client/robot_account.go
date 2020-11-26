package client

import (
	"fmt"
	"strings"

	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func RobotBody(d *schema.ResourceData, projectid string) models.RobotBody {
	resource := strings.Replace(projectid, "s", "", +1)

	body := models.RobotBody{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
	_, existingAttributeOk := d.GetOk("action")
	_, newAttributeOk := d.GetOk("actions")

	robotAccess := models.RobotBodyAccess{}

	if !existingAttributeOk && !newAttributeOk {
		fmt.Errorf("one of action or actions must be configured")

	} else if existingAttributeOk {
		robotAccess.Action = d.Get("action").(string)
		robotAccess.Resource = resource + "/repository"
		body.Access = append(body.Access, robotAccess)

	} else if newAttributeOk {
		access := d.Get("actions").([]interface{})
		for _, v := range access {

			switch v.(string) {
			case "push", "pull":
				robotAccess.Action = v.(string)
				robotAccess.Resource = resource + "/repository"
			case "read":
				robotAccess.Action = v.(string)
				robotAccess.Resource = resource + "/helm-chart"
			case "create":
				robotAccess.Action = v.(string)
				robotAccess.Resource = resource + "/helm-chart-version"
			}
			body.Access = append(body.Access, robotAccess)
		}
	}

	return body
}
