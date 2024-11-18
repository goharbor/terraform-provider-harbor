package provider

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataProjectMemberUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataProjectMemberUsersRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_member_users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataProjectMemberUsersRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	projectId := checkProjectid(d.Get("project_id").(string))
	path := projectId + "/members"

	page := 1
	projectMemberUsersData := make([]map[string]interface{}, 0)
	for {
		resp, _, _, err := apiClient.SendRequest("GET", path+"?page="+strconv.Itoa(page), nil, 200)
		if err != nil {
			return err
		}

		var jsonData []models.ProjectMembersBodyResponses
		err = json.Unmarshal([]byte(resp), &jsonData)
		if err != nil {
			return fmt.Errorf("unable to retrieve Harbor project member users data: %s", err)
		}

		// If there is no data on the current page, we have reached the last page
		if len(jsonData) == 0 {
			break
		}

		for _, v := range jsonData {
			if v.EntityType == "u" {

				projectMemberUserData := map[string]interface{}{
					"id":         path + "/" + strconv.Itoa(v.ID),
					"project_id": checkProjectid(strconv.Itoa(v.ProjectID)),
					"user_name":  v.EntityName,
					"role":       client.RoleTypeNumber(v.RoleID),
				}

				projectMemberUsersData = append(projectMemberUsersData, projectMemberUserData)
			}
		}

		page++
	}
	d.SetId("harbor-project-member-users")
	d.Set("project_member_users", projectMemberUsersData)

	return nil
}
