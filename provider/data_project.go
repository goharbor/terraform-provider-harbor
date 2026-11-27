package provider

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataProject() *schema.Resource {
	return &schema.Resource{
		Read: dataProjectRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
			"vulnerability_scanner": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// lookupPageSize is the page size used when a data source looks a resource up
// by name. Harbor matches the name query as a substring and returns only 10
// results per page by default, so the exact name can be on a later page.
const lookupPageSize = 100

func dataProjectRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	name := d.Get("name").(string)

	for page := 1; ; page++ {
		projectPath := fmt.Sprintf("%s?name=%s&page=%d&page_size=%d", models.PathProjects, name, page, lookupPageSize)

		resp, _, _, err := apiClient.SendRequest("GET", projectPath, nil, 200)
		if err != nil {
			return err
		}

		var jsonData []models.ProjectsBodyResponses
		err = json.Unmarshal([]byte(resp), &jsonData)
		if err != nil {
			return fmt.Errorf("unable to find the project named: %s", name)
		}

		for _, v := range jsonData {

			if v.Name == name {
				id := models.PathProjects + "/" + strconv.Itoa(v.ProjectID)
				public := getboolfromstring(v.Metadata.Public)
				autoScan := getboolfromstring(v.Metadata.AutoScan)
				project_type := getProjectType(v)

				d.SetId(id)
				d.Set("project_id", v.ProjectID)
				d.Set("name", v.Name)
				d.Set("public", public)
				d.Set("vulnerability_scanning", autoScan)
				d.Set("type", project_type)

				scannerName, err := apiClient.GetProjectScanner(id)
				if err != nil {
					return err
				}
				d.Set("vulnerability_scanner", scannerName)
				return nil
			}
		}

		if len(jsonData) < lookupPageSize {
			return nil
		}
	}
}

func getProjectType(project models.ProjectsBodyResponses) string {
	if project.RegistryID != nil {
		return "ProxyCache"
	}
	return "Project"
}

func getboolfromstring(stringbool string) bool {
	boolbool, err := strconv.ParseBool(stringbool)
	if err != nil {
		return false
	}
	return boolbool
}
