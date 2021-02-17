package provider

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRentention() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"schedule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 15,
				MinItems: 1,
				Elem: &schema.Resource{
					// Schema: retentionPolicyRuleFields(),
					Schema: map[string]*schema.Schema{
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"n_days_since_last_pull": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  nil,
							// ConflictsWith: []string{"n_days_since_last_push"},
						},
						"n_days_since_last_push": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  nil,
							// ConflictsWith: []string{"n_days_since_last_pull"},
						},
						"most_recently_pulled": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"most_recently_pushed": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"always_retain": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"repo_matching": {
							Type:     schema.TypeString,
							Optional: true,
							// ConflictsWith: []string{".repo_excluding"},
						},
						"repo_excluding": {
							Type:     schema.TypeString,
							Optional: true,
							// ConflictsWith: []string{".repo_matching"},
						},
						"tag_matching": {
							Type:     schema.TypeString,
							Optional: true,
							// ConflictsWith: []string{".tag_excluding"},
						},
						"tag_excluding": {
							Type:     schema.TypeString,
							Optional: true,
							// ConflictsWith: []string{".tag_matching"},
						},
						"untagged_artifacts": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
		},
		Create: resourceRententionCreate,
		Read:   resourceRententionRead,
		Update: resourceRententionUpdate,
		Delete: resourceRententionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceRententionCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.GetRententionBody(d)
	id := ""

	_, headers, err := apiClient.SendRequest("POST", models.PathRetentions, body, 201)
	if err != nil {
		project_id := strconv.Itoa(body.Scope.Ref)
		resp, _, err := apiClient.SendRequest("GET", models.PathProjects+"/"+project_id, nil, 200)

		var jsonData models.ProjectsBodyResponses
		err = json.Unmarshal([]byte(resp), &jsonData)

		if err != nil {
			return err
		}
		_, headers, err = apiClient.SendRequest("PUT", models.PathRetentions+"/"+jsonData.Metadata.RetentionId, body, 200)
		if err != nil {
			return err
		}
		id = models.PathRetentions + "/" + jsonData.Metadata.RetentionId
	} else {
		id, err = client.GetID(headers)
	}

	d.SetId(id)
	return resourceRententionRead(d, m)
}

func resourceRententionRead(d *schema.ResourceData, m interface{}) error {
	// apiClient := m.(*client.Client)

	// resp, _, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	// if err != nil {
	// 	return err
	// }

	// var jsonData models.Labels
	// err = json.Unmarshal([]byte(resp), &jsonData)
	// if err != nil {
	// 	return fmt.Errorf("Resource not found %s", d.Id())
	// }

	// d.Set("name", jsonData.Name)
	// d.Set("description", jsonData.Description)
	// d.Set("color", jsonData.Color)
	// d.Set("scope", jsonData.Scope)

	return nil
}

func resourceRententionUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.GetRententionBody(d)

	_, _, err := apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		return err
	}

	return resourceRententionRead(d, m)
}

func resourceRententionDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	scope := d.Get("scope").(string)
	project_id, err := strconv.Atoi(strings.ReplaceAll(scope, "/projects/", ""))

	retention := d.Id()
	retention_id, err := strconv.Atoi(strings.ReplaceAll(retention, "/retentions/", ""))

	body := models.Retention{
		Algorithm: "or",
		Scope: models.Scope{
			Level: "project",
			Ref:   project_id,
		},
		Trigger: models.Trigger{
			Kind: "Schedule",
			Settings: models.Settings{
				Cron: "",
			},
		},
		Rules: []models.Rules{},
		Id:    retention_id,
	}

	_, _, err = apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		return err
	}
	return nil
}
