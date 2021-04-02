package client

import (
	"log"
	"regexp"
	"strconv"

	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetRententionBody(d *schema.ResourceData) *models.Retention {
	scope := d.Get("scope").(string)
	re := regexp.MustCompile(`(?m)[0-9]+`)

	id, _ := strconv.Atoi(re.FindString(scope))

	body := models.Retention{
		Algorithm: "or",
		Scope: models.Scope{
			Level: "project",
			Ref:   id,
		},
	}

	body.Trigger = models.Trigger{
		Kind: "Schedule",
	}

	schedule := d.Get("schedule").(string)
	if schedule != "" {
		_, cronstr := GetSchedule(schedule)
		body.Trigger.Settings = models.Settings{
			Cron: cronstr,
		}
	} else {
		body.Trigger.Settings = models.Settings{
			Cron: "",
		}
	}

	body.Rules = expandRententionRules(d)
	log.Printf("[DEBUG] %+v\n ", body)

	return &body
}

func expandRententionRules(d *schema.ResourceData) []models.Rules {
	RetentionRules := d.Get("rule").([]interface{})
	var rules []models.Rules

	for _, v := range RetentionRules {
		i := v.(map[string]interface{})

		rule := models.Rules{
			Disabled: i["disabled"].(bool),
			Action:   "retain",
		}

		if i["n_days_since_last_pull"].(int) > 0 {
			rule.Params.NDaysSinceLastPull = i["n_days_since_last_pull"].(int)
			rule.Template = "nDaysSinceLastPull"
		} else if i["n_days_since_last_push"].(int) > 0 {
			rule.Params.NDaysSinceLastPush = i["n_days_since_last_push"].(int)
			rule.Template = "nDaysSinceLastPush"
		} else if i["most_recently_pulled"].(int) > 0 {
			rule.Params.LatestPulledN = i["most_recently_pulled"].(int)
			rule.Template = "latestPulledN"
		} else if i["most_recently_pushed"].(int) > 0 {
			rule.Params.LatestPushedK = i["most_recently_pushed"].(int)
			rule.Template = "latestPushedK"
		} else if i["always_retain"] == true {
			rule.Template = "always"
		}

		tag := models.TagSelectors{
			Kind: "doublestar",
		}

		if i["tag_matching"].(string) != "" {
			tag.Decoration = "matches"
			tag.Pattern = i["tag_matching"].(string)
			tag.Extras = "{\"untagged\":" + strconv.FormatBool(i["untagged_artifacts"].(bool)) + "}"
			log.Printf("[DEBUG] %s\n ", "tag matching")
		}

		if i["tag_excluding"].(string) != "" {
			tag.Decoration = "excludes"
			tag.Pattern = i["tag_excluding"].(string)
			tag.Extras = "{\"untagged\":" + strconv.FormatBool(i["untagged_artifacts"].(bool)) + "}"
			log.Printf("[DEBUG] %s\n ", "tag excluding")
		}

		scopeSelectorsRepository := models.ScopeSelectors{
			Repository: []models.Repository{},
		}

		repo := models.Repository{
			Kind: "doublestar",
		}

		if i["repo_matching"].(string) != "" {
			repo.Decoration = "repoMatches"
			repo.Pattern = i["repo_matching"].(string)
		}
		if i["repo_excluding"].(string) != "" {
			repo.Decoration = "repoMExcludes"
			repo.Pattern = i["repo_excluding"].(string)
		}

		scopeSelectorsRepository.Repository = append(scopeSelectorsRepository.Repository, repo)

		rule.ScopeSelectors = scopeSelectorsRepository
		rule.TagSelectors = append(rule.TagSelectors, tag)
		rules = append(rules, rule)
	}
	return rules
}
