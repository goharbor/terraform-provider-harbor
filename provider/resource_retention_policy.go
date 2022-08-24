package provider

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"strconv"
	"strings"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRetention() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"schedule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
		CreateContext: resourceRetentionCreate,
		ReadContext:   resourceRetentionRead,
		UpdateContext: resourceRetentionUpdate,
		DeleteContext: resourceRetentionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceRetentionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	body := client.GetRententionBody(d)
	id := ""

	_, headers, _, err := apiClient.SendRequest(ctx, "POST", models.PathRetentions, body, 201)
	if err != nil {
		project_id := strconv.Itoa(body.Scope.Ref)
		resp, _, _, err := apiClient.SendRequest(ctx, "GET", models.PathProjects+"/"+project_id, nil, 200)

		var jsonData models.ProjectsBodyResponses
		err = json.Unmarshal([]byte(resp), &jsonData)

		if err != nil {
			return diag.FromErr(err)
		}
		_, headers, _, err = apiClient.SendRequest(ctx, "PUT", models.PathRetentions+"/"+jsonData.Metadata.RetentionId, body, 200)
		if err != nil {
			return diag.FromErr(err)
		}
		id = models.PathRetentions + "/" + jsonData.Metadata.RetentionId
	} else {
		id, err = client.GetID(headers)
	}

	d.SetId(id)
	return resourceRetentionRead(ctx, d, m)
}

func resourceRetentionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	log.Printf("[DEBUG] Id: %+v\n", d.Id())
	resp, _, respCode, err := apiClient.SendRequest(ctx, "GET", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Retention policy %q was not found - removing from state!", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("making read request on retention policy %s : %+v", d.Id(), err)
	}

	var retentionModel models.Retention
	err = json.Unmarshal([]byte(resp), &retentionModel)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("scope", resolveScope(retentionModel)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("schedule", resolveSchedule(retentionModel)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rule", resolveRules(retentionModel)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRetentionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	body := client.GetRententionBody(d)

	_, _, _, err := apiClient.SendRequest(ctx, "PUT", d.Id(), body, 200)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRetentionRead(ctx, d, m)
}

func resourceRetentionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	_, _, _, err = apiClient.SendRequest(ctx, "PUT", d.Id(), body, 200)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resolveScope(model models.Retention) interface{} {
	return models.PathProjects + "/" + strconv.Itoa(model.Scope.Ref)
}

func resolveSchedule(model models.Retention) string {
	fmt, _ := client.GetSchedule(model.Trigger.Settings.Cron)
	if fmt == "Custom" {
		return model.Trigger.Settings.Cron
	}
	return fmt
}

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
