package provider

import (
	"encoding/json"
	"log"

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

	_, headers, err := apiClient.SendRequest("POST", path, body, 201)
	if err != nil {
		return err
	}

	id, err = client.GetID(headers)
	d.SetId(id)
	return resourceImmutableTagRuleRead(d, m)
}

func resourceImmutableTagRuleRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	log.Printf("[DEBUG] Id: %+v\n", d.Id())
	resp, _, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	if err != nil {
		return err
	}

	var immutableTagRuleModel models.ImmutableTagRule
	err = json.Unmarshal([]byte(resp), &immutableTagRuleModel)
	if err != nil {
		return err
	}

	// GET works
	// log.Printf("[DEBUG] %+v\n", resp)
	// log.Printf("[DEBUG] %+v\n", jsonData)

	//TODO
	//if err := d.Set("rule", resolveRules(retentionModel)); err != nil {
	//	return err
	//}

	return nil
}

func resourceImmutableTagRuleUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.GetImmutableTagRuleBody(d)

	_, _, err := apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		return err
	}

	return resourceImmutableTagRuleRead(d, m)
}

//TODO
func resourceImmutableTagRuleDelete(d *schema.ResourceData, m interface{}) error {
	/*
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
	*/
	return nil
}

//TODO
/*
func resolveRules(model models.Retention) []interface{} {
	modelRules := &model.Rules
	if modelRules != nil {
		flatRules := make([]interface{}, len(*modelRules), len(*modelRules))

		for i, modelRule := range *modelRules {
			flatRule := make(map[string]interface{})

			flatRule["disabled"] = modelRule.Disabled

			switch modelRule.Template {
			case "always":
				flatRule["always_retain"] = true
			case "latestPulledN":
				flatRule["most_recently_pulled"] = modelRule.Params.LatestPulledN
			case "latestPushedK":
				flatRule["most_recently_pushed"] = modelRule.Params.LatestPushedK
			case "nDaysSinceLastPull":
				flatRule["n_days_since_last_pull"] = modelRule.Params.NDaysSinceLastPull
			case "nDaysSinceLastPush":
				flatRule["n_days_since_last_push"] = modelRule.Params.NDaysSinceLastPush
			}

			switch modelRule.TagSelectors[0].Decoration {
			case "matches":
				flatRule["tag_matching"] = modelRule.TagSelectors[0].Pattern
			case "excludes":
				flatRule["tag_excluding"] = modelRule.TagSelectors[0].Pattern
			}

			switch modelRule.ScopeSelectors.Repository[0].Decoration {
			case "repoMatches":
				flatRule["repo_matching"] = modelRule.ScopeSelectors.Repository[0].Pattern
			case "repoExcludes":
				flatRule["repo_excluding"] = modelRule.ScopeSelectors.Repository[0].Pattern
			}

			flatRule["untagged_artifacts"] = strings.Contains(modelRule.TagSelectors[0].Extras, "true")

			flatRules[i] = flatRule
		}

		return flatRules
	}

	return make([]interface{}, 0)
}
*/
