package provider

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataProjectsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"vulnerability_scanning": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"public": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"vulnerability_scanning": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataProjectsRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	name := d.Get("name").(string)
	projectType := d.Get("type").(string)
	public := d.Get("public").(bool)
	vulnScanning := d.Get("vulnerability_scanning").(bool)

	page := 1
	projectsData := make([]map[string]interface{}, 0)
	for {
		resp, _, _, err := apiClient.SendRequest("GET", models.PathProjects+"?page="+strconv.Itoa(page), nil, 200)
		if err != nil {
			return err
		}

		var jsonData []models.ProjectsBodyResponses
		err = json.Unmarshal([]byte(resp), &jsonData)
		if err != nil {
			return fmt.Errorf("unable to retrieve Harbor projects data: %s", err)
		}

		// If there is no data on the current page, we have reached the last page
		if len(jsonData) == 0 {
			break
		}

		for _, v := range jsonData {
			if (name == "" || matchRegex(name, v.Name)) &&
				(projectType == "" || projectType == getProjectType(v)) &&
				(!public || public == getboolfromstring(v.Metadata.Public)) &&
				(!vulnScanning || vulnScanning == getboolfromstring(v.Metadata.AutoScan)) {

				projectData := map[string]interface{}{
					"name":                   v.Name,
					"project_id":             v.ProjectID,
					"public":                 getboolfromstring(v.Metadata.Public),
					"vulnerability_scanning": getboolfromstring(v.Metadata.AutoScan),
				}

				projectData["type"] = getProjectType(v)

				projectsData = append(projectsData, projectData)
			}
		}

		page++
	}
	d.SetId("harbor-projects")
	d.Set("projects", projectsData)

	return nil
}

func matchRegex(pattern, value string) bool {
	match, err := regexp.MatchString(pattern, value)
	return err == nil && match
}
