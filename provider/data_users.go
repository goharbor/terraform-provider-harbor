package provider

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataUsersRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"full_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"admin": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataUsersRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	username := d.Get("username").(string)
	email := d.Get("email").(string)

	usersQueryPath := []string{}
	if username != "" {
		usersQueryPath = append(usersQueryPath, "username="+username)
	}
	if email != "" {
		usersQueryPath = append(usersQueryPath, "email="+email)
	}

	page := 1
	usersData := make([]map[string]interface{}, 0)
	for {
		usersPath := models.PathUsers + "?page=" + strconv.Itoa(page)
		if len(usersQueryPath) > 0 {
			usersPath += "&q=" + strings.Join(usersQueryPath, ",")
		}

		resp, _, _, err := apiClient.SendRequest("GET", usersPath, nil, 200)
		if err != nil {
			return err
		}

		var jsonData []models.UserBody
		err = json.Unmarshal([]byte(resp), &jsonData)
		if err != nil {
			return fmt.Errorf("unable to retrieve Harbor users data: %s", err)
		}

		// If there is no data on the current page, we have reached the last page
		if len(jsonData) == 0 {
			break
		}

		for _, v := range jsonData {
			id := models.PathUsers + "/" + strconv.Itoa(v.UserID)

			userData := map[string]interface{}{
				"id":        id,
				"username":  v.Username,
				"full_name": v.Realname,
				"email":     v.Email,
				"admin":     v.SysadminFlag,
				"comment":   v.Comment,
			}

			usersData = append(usersData, userData)
		}

		page++
	}
	d.SetId("harbor-users")
	d.Set("users", usersData)

	return nil
}
