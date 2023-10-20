package provider

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataGroupsRead,
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ldap_group_dn": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ldap_group_dn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataGroupsRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	groupName := d.Get("group_name").(string)
	ldapGroupDN := d.Get("ldap_group_dn").(string)

	page := 1
	userGroupsData := make([]map[string]interface{}, 0)
	for {
		resp, _, _, err := apiClient.SendRequest("GET", models.PathGroups+"?page="+strconv.Itoa(page), nil, 200)
		if err != nil {
			return err
		}

		var jsonData []models.GroupBody
		err = json.Unmarshal([]byte(resp), &jsonData)
		if err != nil {
			return fmt.Errorf("unable to retrieve Harbor user groups data: %s", err)
		}

		// If there is no data on the current page, we have reached the last page
		if len(jsonData) == 0 {
			break
		}

		for _, v := range jsonData {
			if (groupName == "" || v.Groupname == groupName) &&
				(ldapGroupDN == "" || v.LdapGroupDn == ldapGroupDN) {

				userGroupData := map[string]interface{}{
					"group_name":    v.Groupname,
					"group_type":    v.GroupType,
					"id":            v.ID,
					"ldap_group_dn": v.LdapGroupDn,
				}

				userGroupsData = append(userGroupsData, userGroupData)
			}
		}

		page++
	}
	d.SetId("harbor-user-groups")
	d.Set("groups", userGroupsData)

	return nil
}
