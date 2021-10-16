package client

import (
	"log"
	"regexp"
	"strconv"

	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetTagImmutableRuleBody(d *schema.ResourceData) *models.ImmutableTagRule {
	scope := d.Get("scope").(string)
	re := regexp.MustCompile(`(?m)[0-9]+`)

	id, _ := strconv.Atoi(re.FindString(scope))

	tags := []models.ImmutableTagRuleTagSelectors{}
	tag := models.ImmutableTagRuleTagSelectors{
		Kind: "doublestar",
	}
	tags = append(tags, tag)

	if d.Get("tag_matching").(string) != "" {
		tag.Decoration = "matches"
		tag.Pattern = d.Get("tag_matching").(string)
	}

	if d.Get("tag_excluding").(string) != "" {
		tag.Decoration = "excludes"
		tag.Pattern = d.Get("tag_excluding").(string)
	}

	scopeSelectorsRepository := models.ScopeSelectors{
		Repository: []models.Repository{},
	}

	repo := models.Repository{
		Kind: "doublestar",
	}

	if d.Get("repo_matching").(string) != "" {
		repo.Decoration = "repoMatches"
		repo.Pattern = d.Get("repo_matching").(string)
	}
	if d.Get("repo_excluding").(string) != "" {
		repo.Decoration = "repoExcludes"
		repo.Pattern = d.Get("repo_excluding").(string)
	}
	scopeSelectorsRepository.Repository = append(scopeSelectorsRepository.Repository, repo)

	body := models.ImmutableTagRule{
		Scope: models.Scope{
			Level: "project",
			Ref:   id,
		},
		Disabled:                     d.Get("disabled").(bool),
		ScopeSelectors:               scopeSelectorsRepository,
		ImmutableTagRuleTagSelectors: tags,
	}

	log.Printf("[DEBUG] %+v\n ", body)

	return &body
}
