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
		},
	}
}

func dataProjectRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	name := d.Get("name").(string)
	projectPath := models.PathProjects + "?name=" + name

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
		}
	}
	return nil
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
