package client

import (
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GroupBody return a json body
func GroupBody(d *schema.ResourceData) models.GroupBody {
	return models.GroupBody{
		Groupname:    d.Get("group_name").(string),
		GroupType:    d.Get("group_type").(int),
	}
}
