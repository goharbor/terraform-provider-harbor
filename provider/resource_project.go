package provider

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
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
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	body := client.ProjectBody(d)

	_, headers, _, err := apiClient.SendRequest(ctx, "POST", models.PathProjects, body, 201)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := client.GetID(headers)
	d.SetId(id)
	return resourceProjectRead(ctx, d, m)
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest(ctx, "GET", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Project %q was not found - removing from state!", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("making read request on project %s : %+v", d.Id(), err)
	}

	var jsonData models.ProjectsBodyResponses
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return diag.FromErr(err)
	}

	autoScan := jsonData.Metadata.AutoScan
	var vuln bool
	if autoScan == "" {
		vuln = false
	} else {
		vuln, err = strconv.ParseBool(autoScan)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	var trust bool
	trustContent := jsonData.Metadata.EnableContentTrust
	if trustContent == "" {
		trust = false
	} else {
		trust, err = strconv.ParseBool(trustContent)
		if err != nil {
			return diag.FromErr(err)
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

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	body := client.ProjectBody(d)

	_, _, _, err := apiClient.SendRequest(ctx, "PUT", d.Id(), body, 200)
	if err != nil {
		return diag.FromErr(err)
	}

	apiClient.UpdateStorageQuota(ctx, d)

	return resourceProjectRead(ctx, d, m)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	forceDestroy := d.Get("force_destroy").(bool)

	if forceDestroy {
		// If force_destroy is set delete all repositories within the project
		// before attempting to delete it.
		projectName := d.Get("name").(string)

		err := apiClient.DeleteProjectRepositories(ctx, projectName)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	_, _, respCode, err := apiClient.SendRequest(ctx, "DELETE", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Resource project %q was not found - already deleted!", d.Id())
			return nil
		}
		return diag.Errorf("making delete request on resource project %s : %+v", d.Id(), err)
	}

	return nil
}
