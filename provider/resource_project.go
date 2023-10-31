package provider

import (
	"encoding/json"
	"fmt"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"registry_id": {
				Type:     schema.TypeInt,
				Computed: true,
				ForceNew: true,
				Optional: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"vulnerability_scanning": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"storage_quota": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"deployment_security": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cve_allowlist": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"enable_content_trust": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enable_content_trust_cosign": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.ProjectBody(d)

	_, headers, _, err := apiClient.SendRequest("POST", models.PathProjects, body, 201)
	if err != nil {
		return err
	}

	id, _ := client.GetID(headers)
	d.SetId(id)
	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	if respCode == 404 && err != nil {
		d.SetId("")
		return nil
	} else if err != nil {
		return err
	}

	var jsonData models.ProjectsBodyResponses
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("resource not found %s", d.Id())
	}

	vuln, err := client.ParseBoolOrDefault(jsonData.Metadata.AutoScan, false)
	if err != nil {
		return err
	}

	trust, err := client.ParseBoolOrDefault(jsonData.Metadata.EnableContentTrust, false)
	if err != nil {
		return err
	}

	trustCosign, err := client.ParseBoolOrDefault(jsonData.Metadata.EnableContentTrustCosign, false)
	if err != nil {
		return err
	}

	public, err := client.ParseBoolOrDefault(jsonData.Metadata.Public, false)
	if err != nil {
		return err
	}

	d.Set("name", jsonData.Name)
	d.Set("project_id", jsonData.ProjectID)
	d.Set("registry_id", jsonData.RegistryID)
	d.Set("public", public)
	d.Set("vulnerability_scanning", vuln)
	d.Set("enable_content_trust", trust)
	d.Set("enable_content_trust_cosign", trustCosign)

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.ProjectBody(d)

	_, _, _, err := apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		return err
	}

	apiClient.UpdateStorageQuota(d)

	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	forceDestroy := d.Get("force_destroy").(bool)

	if forceDestroy {
		// If force_destroy is set delete all repositories within the project
		// before attempting to delete it.
		projectName := d.Get("name").(string)

		err := apiClient.DeleteProjectRepositories(projectName)
		if err != nil {
			return err
		}
	}
	if !forceDestroy {
		projectName := d.Get("name").(string)
		repos, _ := apiClient.GetProjectRepositories(projectName)
		if len(repos) != 0 {
			return fmt.Errorf("project %s is not empty, please set force_delete to TRUE to clean all repositories", projectName)
		}
	}

	_, _, respCode, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if respCode != 404 && err != nil { // We can't delete something that doesn't exist. Hence the 404-check
		return err
	}
	return nil
}
