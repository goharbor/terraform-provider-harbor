package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceImmutableTagRule() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
		},
		Create: resourceImmutableTagRuleCreate,
		Read:   resourceImmutableTagRuleRead,
		Update: resourceImmutableTagRuleUpdate,
		Delete: resourceImmutableTagRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceImmutableTagRuleCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	projectid := checkProjectid(d.Get("project_id").(string))
	path := projectid + models.PathImmutableTagRules

	body := client.GetImmutableTagRuleBody(d)
	id := ""

	_, headers, _, err := apiClient.SendRequest("POST", path, body, 201)
	if err != nil {
		return err
	}

	id, err = client.GetID(headers)
	d.SetId(id)
	return resourceImmutableTagRuleRead(d, m)
}

func resourceImmutableTagRuleRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	lastSlashIndex := strings.LastIndex(d.Id(), "/")
	projectImmutableTagRulePath := d.Id()[0:lastSlashIndex]
	tagId, err := strconv.Atoi(d.Id()[lastSlashIndex+1:])
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Path to immutable tag rules: %+v\n", projectImmutableTagRulePath)

	resp, _, _, err := apiClient.SendRequest("GET", projectImmutableTagRulePath, nil, 200)
	if err != nil {
		return fmt.Errorf("Resource not found %s", projectImmutableTagRulePath)
	}

	var immutableTagRuleModels []models.ImmutableTagRule
	err = json.Unmarshal([]byte(resp), &immutableTagRuleModels)
	if err != nil {
		return err
	}
	for _, rule := range immutableTagRuleModels {
		if rule.Id == tagId {
			log.Printf("[DEBUG] found tag id %d", tagId)
			d.Set("disabled", rule.Disabled)
			d.Set("project_id", strings.ReplaceAll(projectImmutableTagRulePath, models.PathImmutableTagRules, ""))

			switch rule.ImmutableTagRuleTagSelectors[0].Decoration {
			case "matches":
				d.Set("tag_matching", rule.ImmutableTagRuleTagSelectors[0].Pattern)
			case "excludes":
				d.Set("tag_excluding", rule.ImmutableTagRuleTagSelectors[0].Pattern)
			}

			switch rule.ScopeSelectors.Repository[0].Decoration {
			case "repoMatches":
				d.Set("repo_matching", rule.ScopeSelectors.Repository[0].Pattern)
			case "excludes":
				d.Set("repo_excluding", rule.ScopeSelectors.Repository[0].Pattern)
			}

			return nil
		}
	}

	return fmt.Errorf("Resource not found %s", d.Id())
}

func resourceImmutableTagRuleUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.GetImmutableTagRuleBody(d)
	id := d.Id()
	log.Printf("[DEBUG] Update Id: %+v\n", id)
	_, _, _, err := apiClient.SendRequest("PUT", id, body, 200)
	if err != nil {
		return err
	}

	return resourceImmutableTagRuleRead(d, m)
}

func resourceImmutableTagRuleDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	_, _, _, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if err != nil {
		return err
	}
	return nil
}
