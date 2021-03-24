package provider

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/BESTSELLER/terraform-provider-harbor/models"
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
				Optional: true,
				ForceNew: true,
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
		},
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.ProjectBody(d)

	_, headers, err := apiClient.SendRequest("POST", models.PathProjects, body, 201)
	if err != nil {
		return err
	}

	id, err := client.GetID(headers)
	d.SetId(id)
	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	if err != nil {
		return err
	}

	var jsonData models.ProjectsBodyResponses
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("Resource not found %s", d.Id())
	}

	vuln, err := strconv.ParseBool(jsonData.Metadata.AutoScan)
	if err != nil {
		return err
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
	d.Set("public", jsonData.Metadata.Public)
	d.Set("vulnerability_scanning", vuln)
	d.Set("enable_content_trust", trust)

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.ProjectBody(d)

	_, _, err := apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		return err
	}

	if d.HasChange("storage_quota") {
		quotaID := "/quotas/" + strings.Replace(d.Id(), "/projects", "", -1)

		storage := d.Get("storage_quota").(int)
		if storage > 0 {
			storage *= 1073741824 // GB
		}
		quota := models.Hard{
			Storage: int64(storage),
		}
		body := models.StorageQuota{quota}

		respBody, _, statusCode, err := apiClient.SendRequestWithStatusCode("PUT", quotaID, body, 0)

		if err != nil {
			return err
		}

		if statusCode == 404 {
			if storage == -1 {
				// don't fail if quota does not exist and no quota is supposed to be set
				// this is normal when `quota_per_project_enable` is `false` in Harbor
			} else {
				name := d.Get("name").(string)
				return fmt.Errorf("[ERROR] Quota for %s does not exist. Quotas can not be set if the project was created with Harbor configuration quota_per_project_enable=false", name)
			}

		} else if statusCode != 200 {
			return fmt.Errorf("[ERROR] unexpected status code got: %v expected: %v \n %v", statusCode, 200, respBody)
		}
	}

	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	_, _, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if err != nil {
		return err
	}
	return nil
}
