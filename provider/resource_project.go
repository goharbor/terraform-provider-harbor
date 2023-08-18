package provider

import (
	"encoding/json"
	"fmt"
	"strconv"

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
				Type:     schema.TypeString,
				Optional: true,
				Default:  "false",
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

	id, err := client.GetID(headers)
	d.SetId(id)
	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	if respCode == 404 && err != nil {
		d.SetId("")
		return fmt.Errorf("Resource not found %s", d.Id())
	}

	var jsonData models.ProjectsBodyResponses
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("Resource not found %s", d.Id())
	}

	autoScan := jsonData.Metadata.AutoScan
	var vuln bool
	if autoScan == "" {
		vuln = false
	} else {
		vuln, err = strconv.ParseBool(autoScan)
		if err != nil {
			return err
		}
	}

	var trust bool
	trustContent := jsonData.Metadata.EnableContentTrust
	if trustContent == "" {
		trust = false
	} else {
		trust, err = strconv.ParseBool(trustContent)
		if err != nil {
			return err
		}
	}

	d.Set("name", jsonData.Name)
	d.Set("project_id", jsonData.ProjectID)
	d.Set("registry_id", jsonData.RegistryID)
	d.Set("public", jsonData.Metadata.Public)
	d.Set("vulnerability_scanning", vuln)
	d.Set("enable_content_trust", trust)

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
			return fmt.Errorf("Project %s is not empty, please set force_delete to TRUE to clean all repositories", projectName)
		}
	}

	_, _, _, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if err != nil {
		return err
	}
	return nil
}
