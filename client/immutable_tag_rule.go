package client

import (
	"log"

	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetImmutableTagRuleBody(d *schema.ResourceData) *models.ImmutableTagRule {
	tags := []models.ImmutableTagRuleTagSelectors{}
	tag := models.ImmutableTagRuleTagSelectors{
		Kind: "doublestar",
	}

	if d.Get("tag_matching").(string) != "" {
		tag.Decoration = "matches"
		tag.Pattern = d.Get("tag_matching").(string)
	}

	if d.Get("tag_excluding").(string) != "" {
		tag.Decoration = "excludes"
		tag.Pattern = d.Get("tag_excluding").(string)
	}
	tags = append(tags, tag)

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
		Disabled:                     d.Get("disabled").(bool),
		ScopeSelectors:               scopeSelectorsRepository,
		ImmutableTagRuleTagSelectors: tags,
		Action:                       "immutable",
		Template:                     "immutable_template",
	}

	log.Printf("[DEBUG] %+v\n ", body)

	return &body
}
