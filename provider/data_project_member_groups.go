package provider

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataProjectMemberGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataProjectMemberGroupsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_member_groups": {
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
						"group_name": {
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

func dataProjectMemberGroupsRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	projectId := checkProjectid(d.Get("project_id").(string))
	path := projectId + "/members"

	page := 1
	projectMemberGroupsData := make([]map[string]interface{}, 0)
	for {
		resp, _, _, err := apiClient.SendRequest("GET", path+"?page="+strconv.Itoa(page), nil, 200)
		if err != nil {
			return err
		}

		var jsonData []models.ProjectMembersBodyResponses
		err = json.Unmarshal([]byte(resp), &jsonData)
		if err != nil {
			return fmt.Errorf("unable to retrieve Harbor project member groups data: %s", err)
		}

		// If there is no data on the current page, we have reached the last page
		if len(jsonData) == 0 {
			break
		}

		for _, v := range jsonData {
			if v.EntityType == "g" {

				projectMemberGroupData := map[string]interface{}{
					"id":         path + "/" + strconv.Itoa(v.ID),
					"project_id": checkProjectid(strconv.Itoa(v.ProjectID)),
					"group_name": v.EntityName,
					"role":       client.RoleTypeNumber(v.RoleID),
				}

				projectMemberGroupsData = append(projectMemberGroupsData, projectMemberGroupData)
			}
		}

		page++
	}
	d.SetId("harbor-project-member-groups")
	d.Set("project_member_groups", projectMemberGroupsData)

	return nil
}
