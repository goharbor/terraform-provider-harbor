package provider

import (
	"encoding/json"
	"fmt"
	"strconv"

	"bitbucket.org/bestsellerit/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var pathProjects = "/projects"

type projectsRequest struct {
	CountLimit   int    `json:"count_limit,omitempty"`
	ProjectName  string `json:"project_name,omitempty"`
	CveWhitelist struct {
		Items []struct {
			CveID string `json:"cve_id,omitempty"`
		} `json:"items,omitempty"`
		ProjectID int `json:"project_id,omitempty"`
		ID        int `json:"id,omitempty"`
		ExpiresAt int `json:"expires_at,omitempty"`
	} `json:"cve_whitelist,omitempty"`
	StorageLimit int `json:"storage_limit,omitempty"`
	Metadata     struct {
		EnableContentTrust   string `json:"enable_content_trust,omitempty"`
		AutoScan             string `json:"auto_scan,omitempty"`
		Severity             string `json:"severity,omitempty"`
		ReuseSysCveWhitelist string `json:"reuse_sys_cve_whitelist,omitempty"`
		Public               string `json:"public,omitempty"`
		PreventVul           string `json:"prevent_vul,omitempty"`
	} `json:"metadata,omitempty"`
}

type projectsResponses struct {
	UpdateTime         string `json:"update_time"`
	OwnerName          string `json:"owner_name"`
	Name               string `json:"name"`
	Deleted            bool   `json:"deleted"`
	OwnerID            int    `json:"owner_id"`
	RepoCount          int    `json:"repo_count"`
	CreationTime       string `json:"creation_time"`
	Togglable          bool   `json:"togglable"`
	ProjectID          int    `json:"project_id"`
	CurrentUserRoleID  int    `json:"current_user_role_id"`
	CurrentUserRoleIds []int  `json:"current_user_role_ids"`
	ChartCount         int    `json:"chart_count"`
	CveWhitelist       struct {
		Items []struct {
			CveID string `json:"cve_id"`
		} `json:"items"`
		ProjectID int `json:"project_id"`
		ID        int `json:"id"`
		ExpiresAt int `json:"expires_at"`
	} `json:"cve_whitelist"`
	Metadata struct {
		EnableContentTrust   string `json:"enable_content_trust"`
		AutoScan             string `json:"auto_scan"`
		Severity             string `json:"severity"`
		ReuseSysCveWhitelist string `json:"reuse_sys_cve_whitelist"`
		Public               string `json:"public"`
		PreventVul           string `json:"prevent_vul"`
	} `json:"metadata"`
}

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  false,
			},
			"vulnerability_scanning": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  true,
			},
		},
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	apiClient, body := projectAPIClientRequest(d, m)

	_, err := apiClient.SendRequest("POST", pathProjects, body, 201)
	if err != nil {
		return err
	}

	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	ProjectName := d.Get("name").(string)

	resp, err := apiClient.SendRequest("GET", pathProjects+"?name="+ProjectName, nil, 200)
	if err != nil {
		return err
	}

	var jsonData []projectsResponses
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to unmarchal: %s", err)
	}

	if len(jsonData) < 1 {
		return fmt.Errorf("[ERROR] JsonData is empty")
	}

	for _, v := range jsonData {
		if v.Name == d.Get("name").(string) {
			d.SetId(strconv.Itoa(v.ProjectID))

		}

	}

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient, body := projectAPIClientRequest(d, m)

	_, err := apiClient.SendRequest("PUT", pathProjects, body, 200)
	if err != nil {
		return err
	}

	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	id := d.Id()

	apiClient.SendRequest("DELETE", pathProjects+"/"+id, nil, 200)
	return nil
}

func projectAPIClientRequest(d *schema.ResourceData, m interface{}) (*client.Client, projectsRequest) {
	apiClient := m.(*client.Client)

	body := projectsRequest{
		ProjectName: d.Get("name").(string),
	}
	body.Metadata.AutoScan = d.Get("vulnerability_scanning").(string)
	body.Metadata.Public = d.Get("public").(string)

	return apiClient, body
}
