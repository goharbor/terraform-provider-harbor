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

// CreateUser will create an interal user
// func (client *Client) CreateUser(body *models.UserBody) (*models.UserBody, error) {
// 	path := models.PathUsers

// 	resp, err := client.SendRequest("POST", path, body, 201)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var jsonData models.UserBody
// 	err = json.Unmarshal([]byte(resp), &jsonData)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &jsonData, nil
// }

// UpdateUser will update an Interal user
// func (client *Client) UpdateUser(id string, d *schema.ResourceData) (json string, error) {

// 	body := models.UserBody{
// 		Password:     d.Get("password").(string),
// 		Email:        d.Get("email").(string),
// 		Realname:     d.Get("full_name").(string),
// 	}

// 	_, err := client.SendRequest("PUT", id, updatebody, 200)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return nil, nil
// }
