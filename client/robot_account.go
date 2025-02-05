package client

import (
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func RobotBody(d *schema.ResourceData) models.RobotBody {
	body := models.RobotBody{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Disable:     d.Get("disable").(bool),
		Duration:    d.Get("duration").(int),
		Level:       d.Get("level").(string),
	}

	permissions := d.Get("permissions").(*schema.Set).List()
	for _, p := range permissions {

		permission := models.RobotBodyPermission{
			Kind:      p.(map[string]interface{})["kind"].(string),
			Namespace: p.(map[string]interface{})["namespace"].(string),
		}

		for _, a := range p.(map[string]interface{})["access"].(*schema.Set).List() {
			access := models.RobotBodyAccess{
				Action:   a.(map[string]interface{})["action"].(string),
				Resource: a.(map[string]interface{})["resource"].(string),
			}
			if a.(map[string]interface{})["effect"] != "" {
				access.Effect = a.(map[string]interface{})["effect"].(string)
			}

			permission.Access = append(permission.Access, access)
		}

		body.Permissions = append(body.Permissions, permission)
	}

	return body
}
