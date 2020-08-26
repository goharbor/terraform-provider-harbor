package client

import (
	// "github.com/BESTSELLER/terraform-provider-habor/client"

	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// UserBody return a json body
func UserBody(d *schema.ResourceData) models.UserBody {
	return models.UserBody{
		Username:     d.Get("username").(string),
		Password:     d.Get("password").(string),
		SysadminFlag: d.Get("admin").(bool),
		Email:        d.Get("email").(string),
		Realname:     d.Get("full_name").(string),
		Newpassword:  d.Get("password").(string),
	}
}
