package provider

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"strconv"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataProjectRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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

func dataProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	name := d.Get("name").(string)
	projectPath := models.PathProjects + "?name=" + name

	resp, _, _, err := apiClient.SendRequest(ctx, "GET", projectPath, nil, 200)
	if err != nil {
		return diag.FromErr(err)
	}

	var jsonData []models.ProjectsBodyResponses
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range jsonData {

		if v.Name == name {
			id := models.PathProjects + "/" + strconv.Itoa(v.ProjectID)
			public, err := strconv.ParseBool(v.Metadata.Public)
			if err != nil {
				return diag.FromErr(err)
			}

			var autoScan bool
			scan := v.Metadata.AutoScan
			if scan == "" {
				autoScan = false
			} else {
				autoScan, err = strconv.ParseBool(scan)
				if err != nil {
					return diag.FromErr(err)
				}
			}

			d.SetId(id)
			d.Set("project_id", v.ProjectID)
			d.Set("name", v.Name)
			d.Set("public", public)
			d.Set("vulnerability_scanning", autoScan)
		}
	}
	return nil
}
