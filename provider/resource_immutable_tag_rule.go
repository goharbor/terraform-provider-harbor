package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
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

	_, headers, respCode, err := apiClient.SendRequest("POST", path, body, 201)
	if respCode == 409 {
		log.Printf("[DEBUG] resource already exists %s", path)
		return err
	} else if err != nil {
		return err
	}

	id, _ = client.GetID(headers)
	d.SetId(id)

	if d.Get("disabled").(bool) {
		// if the rule is disabled, we need to do a second request to disable it as Harbor API doesn't allow to create a disabled rule
		body := client.GetImmutableTagRuleBody(d)
		_, _, _, err := apiClient.SendRequest("PUT", id, body, 200)
		if err != nil {
			return err
		}
	}

	return resourceImmutableTagRuleRead(d, m)
}

func resourceImmutableTagRuleRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	lastSlashIndex := strings.LastIndex(d.Id(), "/")
	projectImmutableTagRulePath := d.Id()[0:lastSlashIndex]
	immutableTagRuleId, err := strconv.Atoi(d.Id()[lastSlashIndex+1:])
	if err != nil {
		return err
	}

	var immutableTagRuleModels []models.ImmutableTagRule
	page := 1
	pageSize := 15
	for {
		resp, _, respCode, err := apiClient.SendRequest("GET", fmt.Sprintf("%s?page=%d&page_size=%d", projectImmutableTagRulePath, page, pageSize), nil, 200)
		if respCode == 404 && err != nil {
			d.SetId("")
			return nil
		} else if err != nil {
			return fmt.Errorf("resource not found %s / %s / %d", projectImmutableTagRulePath, err, respCode)
		}
		var pageModels []models.ImmutableTagRule
		err = json.Unmarshal([]byte(resp), &pageModels)
		if err != nil {
			return err
		}
		immutableTagRuleModels = append(immutableTagRuleModels, pageModels...)
		if len(pageModels) < pageSize {
			break
		}
		page++
	}

	for _, rule := range immutableTagRuleModels {
		if rule.Id == immutableTagRuleId {
			log.Printf("[DEBUG] found tag id %d", immutableTagRuleId)
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

	log.Printf("[DEBUG] resource not found %s", d.Id())
	d.SetId("")
	return nil
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
