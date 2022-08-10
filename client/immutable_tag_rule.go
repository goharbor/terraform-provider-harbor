package client

import (
	"log"
	"strconv"
	"strings"

	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetImmutableTagRuleBody(d *schema.ResourceData) *models.ImmutableTagRule {
	tagId := 0
	if d.Id() != "" {
		lastSlashIndex := strings.LastIndex(d.Id(), "/")
		tagId, _ = strconv.Atoi(d.Id()[lastSlashIndex+1:])
	}

	tags := []models.ImmutableTagRuleTagSelectors{}
	tag := models.ImmutableTagRuleTagSelectors{
		Kind: "doublestar",
	}

	if d.Get("tag_matching").(string) != "" {
		tag.Decoration = "matches"
		tag.Pattern = d.Get("tag_matching").(string)
		d.Set("tag_excluding", nil)
	}
	if d.Get("tag_excluding").(string) != "" {
		tag.Decoration = "excludes"
		tag.Pattern = d.Get("tag_excluding").(string)
		d.Set("tag_matching", nil)
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
		d.Set("repo_excluding", nil)
	}
	if d.Get("repo_excluding").(string) != "" {
		repo.Decoration = "repoExcludes"
		repo.Pattern = d.Get("repo_excluding").(string)
		d.Set("repo_matching", nil)
	}
	scopeSelectorsRepository.Repository = append(scopeSelectorsRepository.Repository, repo)

	body := models.ImmutableTagRule{
		Id:                           tagId,
		Disabled:                     d.Get("disabled").(bool),
		ScopeSelectors:               scopeSelectorsRepository,
		ImmutableTagRuleTagSelectors: tags,
		Action:                       "immutable",
		Template:                     "immutable_template",
	}

	log.Printf("[DEBUG] %+v\n ", body)

	return &body
}
