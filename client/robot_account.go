package client

import (
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func RobotBody(d *schema.ResourceData, resource string) models.RobotBody {
	return models.RobotBody{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Access: []models.RobotBodyAccess{
			{
				Action:   d.Get("action").(string),
				Resource: resource,
			},
		},
	}
}
